/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backupscheduled

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/internal/clock"
	"github.com/universityofadelaide/shepherd-operator/internal/events"
	metautils "github.com/universityofadelaide/shepherd-operator/internal/k8s/metadata"
	scheduledutils "github.com/universityofadelaide/shepherd-operator/internal/scheduled"
)

const (
	// ControllerName is used to identify this controller in logs and events.
	ControllerName = "imagescheduled-controller"

	// OwnerKey is used to query Image objects owned by an ImageScheduled object.
	OwnerKey = ".metadata.controller"
)

// Reconciler will manage Backup objects based on a schedule.
type Reconciler struct {
	clock.Clock
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Params   Params
}

// Params used by this controller.
type Params struct {
	// Used to filter Backup objects by a key and value pair.
	FilterByLabelAndValue FilterByLabelAndValue
}

// FilterByLabelAndValue is used to filter Backup objects by a key and value pair.
type FilterByLabelAndValue struct {
	Key   string
	Value string
}

//+kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds/finalizers,verbs=update

// Reconcile will manage Backup objects based on a schedule.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconcile loop")

	scheduled := &extensionv1.BackupScheduled{}

	err := r.Get(context.TODO(), req.NamespacedName, scheduled)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if !metautils.HasLabelWithValue(scheduled.ObjectMeta.Labels, r.Params.FilterByLabelAndValue.Key, r.Params.FilterByLabelAndValue.Value) {
		return reconcile.Result{}, nil
	}

	if scheduled.Spec.Schedule.CronTab == "" {
		return reconcile.Result{}, fmt.Errorf("BackupScheduled doesn't have a schedule.")
	}

	if _, found := scheduled.ObjectMeta.GetLabels()["site"]; !found {
		return reconcile.Result{}, fmt.Errorf("BackupScheduled doesn't have a site label.")
	}

	logger.Info("Querying Backups")

	var backups extensionv1.BackupList

	if err := r.List(ctx, &backups, client.InNamespace(req.Namespace), client.MatchingFields{OwnerKey: req.Name}); err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Enforcing backup retention policies")

	err = r.ExecuteRetentionPolicies(logger, scheduled, backups)
	if err != nil {
		logger.Info(err.Error())
	}

	logger.Info("Filtering Backups")

	active, successful, failed, err := r.SortBackups(scheduled, backups)
	if err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Job count", "active", len(active), "success", len(successful), "failed", len(failed))

	err = r.Status().Update(context.TODO(), scheduled)
	if err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Cleaning up old Backups")

	err = r.Cleanup(scheduled, successful, failed)
	if err != nil {
		return reconcile.Result{}, err
	}

	if scheduled.Spec.Schedule.Suspend != nil && *scheduled.Spec.Schedule.Suspend {
		logger.Info("Scheduling has been suspended. Skipping.")
		return reconcile.Result{}, nil
	}

	return r.ScheduleNextBackup(logger, scheduled, active)
}

// ExecuteRetentionPolicies will delete old Backup objects.
func (r *Reconciler) ExecuteRetentionPolicies(logger logr.Logger, scheduled *extensionv1.BackupScheduled, backups extensionv1.BackupList) error {
	if scheduled.Spec.Retention.MaxNumber == nil || *scheduled.Spec.Retention.MaxNumber < 1 {
		logger.Info("backup retention disabled - skipping")
		return nil
	}

	sort.SliceStable(backups.Items, func(i, j int) bool {
		// Sort by completed date first.
		a := backups.Items[i].Status.CompletionTime
		b := backups.Items[j].Status.CompletionTime
		if a != nil && b != nil {
			return a.After(b.Time)
		}
		// Default to alphabetical.
		return backups.Items[i].Name > backups.Items[j].Name
	})

	remaining := *scheduled.Spec.Retention.MaxNumber

	for _, item := range backups.Items {
		if remaining > 0 {
			logger.Info("keeping backup %s", item.Name)
			remaining--
			continue
		}

		r.Recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventCreate, "Deleting Backup: %s", item.ObjectMeta.Name)

		if err := r.Delete(context.Background(), &item); err != nil {
			return errors.Wrap(err, "failed to delete Backup")
		}
	}

	return nil
}

// SortBackups into active, successful, failed.
func (r *Reconciler) SortBackups(scheduled *extensionv1.BackupScheduled, backups extensionv1.BackupList) ([]*extensionv1.Backup, []*extensionv1.Backup, []*extensionv1.Backup, error) {
	var (
		active       []*extensionv1.Backup
		successful   []*extensionv1.Backup
		failed       []*extensionv1.Backup
		lastSchedule *time.Time
	)

	for i, backup := range backups.Items {
		switch backup.Status.Phase {
		case shpmetav1.PhaseFailed:
			failed = append(failed, &backups.Items[i])
		case shpmetav1.PhaseCompleted:
			successful = append(successful, &backups.Items[i])
		default:
			// If it doesn't have a status assigned yet then the Backup is most likely still be created.
			active = append(active, &backups.Items[i])
		}

		// We'll store the launch time in an annotation, so we'll reconstitute that from
		// the active Backups themselves.
		scheduledTime, err := scheduledutils.GetScheduledTime(backup.Annotations)
		if err != nil {
			return active, successful, failed, err
		}

		if scheduledTime != nil {
			if lastSchedule == nil {
				lastSchedule = scheduledTime
			} else if lastSchedule.Before(*scheduledTime) {
				lastSchedule = scheduledTime
			}
		}
	}

	if lastSchedule != nil {
		scheduled.Status.LastExecutedTime = &metav1.Time{Time: *lastSchedule}
	} else {
		scheduled.Status.LastExecutedTime = nil
	}

	scheduled.Status.Active = nil

	for _, a := range active {
		jobRef, err := reference.GetReference(r.Scheme, a)
		if err != nil {
			continue
		}

		scheduled.Status.Active = append(scheduled.Status.Active, *jobRef)
	}

	return active, successful, failed, nil
}

// Cleanup old successful and failed Backups.
func (r *Reconciler) Cleanup(scheduled *extensionv1.BackupScheduled, successful, failed []*extensionv1.Backup) error {
	if scheduled.Spec.Schedule.FailedHistoryLimit != nil {
		sort.Slice(failed, func(i, j int) bool {
			if failed[i].Status.StartTime == nil {
				return failed[j].Status.StartTime != nil
			}

			return failed[i].Status.StartTime.Before(failed[j].Status.StartTime)
		})
		for i, backup := range failed {
			if int32(i) >= int32(len(failed))-*scheduled.Spec.Schedule.FailedHistoryLimit {
				break
			}

			if err := r.Delete(context.Background(), backup, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old failed job")
			}

			r.Recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting failed Backup: %s", backup.ObjectMeta.Name)
		}
	}

	if scheduled.Spec.Schedule.SuccessfulHistoryLimit != nil {
		sort.Slice(successful, func(i, j int) bool {
			if successful[i].Status.StartTime == nil {
				return successful[j].Status.StartTime != nil
			}

			return successful[i].Status.StartTime.Before(successful[j].Status.StartTime)
		})

		for i, backup := range successful {
			if int32(i) >= int32(len(successful))-*scheduled.Spec.Schedule.SuccessfulHistoryLimit {
				break
			}

			if err := r.Delete(context.Background(), backup, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old successful job")
			}

			r.Recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting successful Backup: %s", backup.ObjectMeta.Name)
		}
	}

	return nil
}

// ScheduleNextBackup checks if a new Backup should be created.
func (r *Reconciler) ScheduleNextBackup(logger logr.Logger, scheduled *extensionv1.BackupScheduled, active []*extensionv1.Backup) (reconcile.Result, error) {
	missedRun, nextRun, err := scheduledutils.GetNextSchedule(scheduled.Spec.Schedule, scheduled.Status, scheduled.ObjectMeta.CreationTimestamp.Time, r.Now())
	if err != nil {
		return reconcile.Result{}, err
	}

	result := reconcile.Result{RequeueAfter: nextRun.Sub(r.Now())}

	if missedRun.IsZero() {
		return result, nil
	}

	logger.Info("ScheduleNextBackup", "missed", missedRun.Unix(), "next", nextRun.Unix())

	// make sure we're not too late to start the run
	tooLate := false

	if scheduled.Spec.Schedule.StartingDeadlineSeconds != nil {
		tooLate = missedRun.Add(time.Duration(*scheduled.Spec.Schedule.StartingDeadlineSeconds) * time.Second).Before(r.Now())
	}

	logger.Info("ScheduleNextBackup", "scheduled.Spec.Schedule.StartingDeadlineSeconds", scheduled.Spec.Schedule.StartingDeadlineSeconds)

	if tooLate {
		return result, nil
	}

	// figure out how to run this job -- concurrency policy might forbid us from running
	// multiple at the same time...
	if scheduled.Spec.Schedule.ConcurrencyPolicy == shpmetav1.ForbidConcurrent && len(active) > 0 {
		return result, nil
	}

	// ...or instruct us to replace existing ones...
	if scheduled.Spec.Schedule.ConcurrencyPolicy == shpmetav1.ReplaceConcurrent {
		for _, activeJob := range active {
			// we don't care if the job was already deleted
			if err := r.Delete(context.Background(), activeJob, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				return reconcile.Result{}, err
			}
		}
	}

	backup, err := buildBackup(scheduled, r.Scheme, missedRun)
	if err != nil {
		return result, errors.Wrap(err, "failed to generate Backup")
	}

	r.Recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventCreate, "Creating Backup: %s", backup.ObjectMeta.Name)

	existing := &extensionv1.Backup{}
	if err = r.Get(context.TODO(), types.NamespacedName{
		Name:      backup.ObjectMeta.Name,
		Namespace: backup.ObjectMeta.Namespace,
	}, existing); err == nil {
		logger.Info("[ScheduleNextBackup] Backup already created", "existing", existing.ObjectMeta.Name, "new", backup.ObjectMeta.Name)
		return result, nil
	}

	if err := r.Create(context.Background(), backup); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create Backup")
	}

	logger.Info("[ScheduleNextBackup] Backup job created", "name", backup.ObjectMeta.Name)

	return result, nil
}

// Helper function to build a backup object.
func buildBackup(scheduled *extensionv1.BackupScheduled, scheme *runtime.Scheme, scheduledTime time.Time) (*extensionv1.Backup, error) {
	backup := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%d", scheduled.Name, scheduledTime.Unix()),
			Namespace: scheduled.ObjectMeta.Namespace,
			Labels:    scheduled.Labels,
		},
		Spec: extensionv1.BackupSpec{
			MySQL:   scheduled.Spec.MySQL,
			Volumes: scheduled.Spec.Volumes,
		},
	}
	if err := controllerutil.SetControllerReference(scheduled, backup, scheme); err != nil {
		return nil, err
	}

	return backup, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Clock == nil {
		r.Clock = clock.New()
	}

	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &extensionv1.Backup{}, OwnerKey, func(rawObj client.Object) []string {
		job := rawObj.(*extensionv1.Backup)

		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != extensionv1.GroupVersion.String() || owner.Kind != "ImageScheduled" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.BackupScheduled{}).
		Owns(&extensionv1.Backup{}).
		Complete(r)
}
