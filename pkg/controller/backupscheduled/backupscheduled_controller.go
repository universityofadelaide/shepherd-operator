package backupscheduled

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ref "k8s.io/client-go/tools/reference"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/clock"
	errutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/controller/errors"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/controller/logger"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/events"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/restic"
	scheduledutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/scheduled"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "backup-scheduled-controller"
	// OwnerKey used to query for child Backups.
	OwnerKey = ".metadata.controller"
)

// Add creates a new BackupScheduled Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackupScheduled{
		Clock:    clock.New(),
		Client:   mgr.GetClient(),
		recorder: mgr.GetRecorder(ControllerName),
		scheme:   mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = mgr.GetFieldIndexer().IndexField(&extensionv1.Backup{}, OwnerKey, func(rawObj runtime.Object) []string {
		backup := rawObj.(*extensionv1.Backup)

		owner := metav1.GetControllerOf(backup)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != extensionv1.SchemeGroupVersion.String() || owner.Kind != "BackupScheduled" {
			return nil
		}

		return []string{owner.Name}
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &extensionv1.BackupScheduled{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileBackupScheduled{}

// ReconcileBackupScheduled reconciles a BackupScheduled object
type ReconcileBackupScheduled struct {
	clock.Clock
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
}

// Reconcile reads that state of the cluster for a BackupScheduled object and makes changes based on the state read
// and what is in the BackupScheduled.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds/finalizers,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileBackupScheduled) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)
	log.Info("Starting reconcile loop")

	scheduled := &extensionv1.BackupScheduled{}
	err := r.Get(context.TODO(), request.NamespacedName, scheduled)
	if err != nil {
		return reconcile.Result{}, errutils.IgnoreNotFound(err)
	}

	if scheduled.Spec.Schedule.CronTab == "" {
		err := errors.New("BackupScheduled doesn't have a schedule.")
		log.Error(err.Error())
		return reconcile.Result{}, err
	}

	if _, found := scheduled.ObjectMeta.GetLabels()["site"]; !found {
		err := errors.New("BackupScheduled doesn't have a site label.")
		log.Error(err.Error())
		return reconcile.Result{}, err
	}

	log.Info("Querying Backups")
	var backups extensionv1.BackupList
	listOptions := client.MatchingField(OwnerKey, request.Name)
	listOptions.Namespace = request.Namespace
	err = r.List(context.Background(), listOptions, &backups)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Filtering Backups")
	active, successful, failed, err := r.SortBackups(scheduled, backups)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Infof("Job count: active=%d success=%d failed=%d", len(active), len(successful), len(failed))

	err = r.Status().Update(context.TODO(), scheduled)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Cleaning up old Backups")
	err = r.Cleanup(scheduled, successful, failed)
	if err != nil {
		return reconcile.Result{}, err
	}

	if scheduled.Spec.Schedule.Suspend != nil && *scheduled.Spec.Schedule.Suspend {
		log.Info("Scheduling has been suspended. Skipping.")
		return reconcile.Result{}, nil
	}

	return r.ScheduleNextBackup(scheduled, active)
}

// SortBackups into active, successful, failed.
func (r *ReconcileBackupScheduled) SortBackups(scheduled *extensionv1.BackupScheduled, backups extensionv1.BackupList) ([]*extensionv1.Backup, []*extensionv1.Backup, []*extensionv1.Backup, error) {
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
		jobRef, err := ref.GetReference(r.scheme, a)
		if err != nil {
			continue
		}

		scheduled.Status.Active = append(scheduled.Status.Active, *jobRef)
	}

	return active, successful, failed, nil
}

// Cleanup old successful and failed Backups.
func (r *ReconcileBackupScheduled) Cleanup(scheduled *extensionv1.BackupScheduled, successful, failed []*extensionv1.Backup) error {
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

			if err := r.Delete(context.Background(), backup, client.PropagationPolicy(metav1.DeletePropagationBackground)); errutils.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old failed job")
			}

			r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting failed Backup: %s", backup.ObjectMeta.Name)
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

			if err := r.Delete(context.Background(), backup, client.PropagationPolicy(metav1.DeletePropagationBackground)); errutils.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old successful job")
			}

			r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting successful Backup: %s", backup.ObjectMeta.Name)
		}
	}

	return nil
}

// ScheduleNextBackup checks if a new Backup should be created.
func (r *ReconcileBackupScheduled) ScheduleNextBackup(scheduled *extensionv1.BackupScheduled, active []*extensionv1.Backup) (reconcile.Result, error) {
	missedRun, nextRun, err := scheduledutils.GetNextSchedule(scheduled.Spec.Schedule, scheduled.Status, scheduled.ObjectMeta.CreationTimestamp.Time, r.Now())
	if err != nil {
		return reconcile.Result{}, err
	}

	result := reconcile.Result{RequeueAfter: nextRun.Sub(r.Now())}

	if missedRun.IsZero() {
		return result, nil
	}

	// make sure we're not too late to start the run
	tooLate := false

	if scheduled.Spec.Schedule.StartingDeadlineSeconds != nil {
		tooLate = missedRun.Add(time.Duration(*scheduled.Spec.Schedule.StartingDeadlineSeconds) * time.Second).Before(r.Now())
	}

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
			if err := r.Delete(context.Background(), activeJob, client.PropagationPolicy(metav1.DeletePropagationBackground)); errutils.IgnoreNotFound(err) != nil {
				return reconcile.Result{}, err
			}
		}
	}

	backup, err := buildBackup(scheduled, r.scheme, missedRun)
	if err != nil {
		return result, errors.Wrap(err, "failed to generate Backup")
	}

	r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventCreate, "Creating Backup: %s", backup.ObjectMeta.Name)

	if err := r.Create(context.Background(), backup); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create Backup")
	}

	return result, nil
}

func buildBackup(scheduled *extensionv1.BackupScheduled, scheme *runtime.Scheme, scheduledTime time.Time) (*extensionv1.Backup, error) {
	backup := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%d", scheduled.Name, scheduledTime.Unix()),
			Namespace: scheduled.ObjectMeta.Namespace,
			Labels:    scheduled.Labels,
			Annotations: map[string]string{
				restic.FriendlyNameAnnotation: scheduledTime.Format(shpmetav1.FriendlyNameFormat),
			},
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
