package scheduled

import (
	"context"
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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	"github.com/skpr/operator/pkg/utils/clock"
	errutils "github.com/skpr/operator/pkg/utils/controller/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/k8s/events"
	scheduledutils "github.com/skpr/operator/pkg/utils/scheduled"
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

// add adds a new Controller to mgr with r as the reconcile.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = mgr.GetFieldIndexer().IndexField(&extensionsv1beta1.Backup{}, OwnerKey, func(rawObj runtime.Object) []string {
		backup := rawObj.(*extensionsv1beta1.Backup)

		owner := metav1.GetControllerOf(backup)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != extensionsv1beta1.SchemeGroupVersion.String() || owner.Kind != "BackupScheduled" {
			return nil
		}

		return []string{owner.Name}
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &extensionsv1beta1.BackupScheduled{}}, &handler.EnqueueRequestForObject{})
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
// Automatically generate RBAC rules to allow the Controller to read and write Backups
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backupscheduleds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backupscheduleds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backups/status,verbs=get;update;patch
func (r *ReconcileBackupScheduled) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)
	log.Info("Starting reconcile loop")

	scheduled := &extensionsv1beta1.BackupScheduled{}
	err := r.Get(context.TODO(), request.NamespacedName, scheduled)
	if err != nil {
		return reconcile.Result{}, errutils.IgnoreNotFound(err)
	}

	log.Info("Querying Backups")
	var backups extensionsv1beta1.BackupList
	err = r.List(context.Background(), client.MatchingField(OwnerKey, request.Name), &backups)
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
func (r *ReconcileBackupScheduled) SortBackups(scheduled *extensionsv1beta1.BackupScheduled, backups extensionsv1beta1.BackupList) ([]*extensionsv1beta1.Backup, []*extensionsv1beta1.Backup, []*extensionsv1beta1.Backup, error) {
	var (
		active       []*extensionsv1beta1.Backup
		successful   []*extensionsv1beta1.Backup
		failed       []*extensionsv1beta1.Backup
		lastSchedule *time.Time
	)

	for i, backup := range backups.Items {
		switch backup.Status.Phase {
		case skprmetav1.PhaseFailed:
			failed = append(failed, &backups.Items[i])
		case skprmetav1.PhaseCompleted:
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
		scheduled.Status.LastScheduleTime = &metav1.Time{Time: *lastSchedule}
	} else {
		scheduled.Status.LastScheduleTime = nil
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
func (r *ReconcileBackupScheduled) Cleanup(scheduled *extensionsv1beta1.BackupScheduled, successful, failed []*extensionsv1beta1.Backup) error {
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
func (r *ReconcileBackupScheduled) ScheduleNextBackup(scheduled *extensionsv1beta1.BackupScheduled, active []*extensionsv1beta1.Backup) (reconcile.Result, error) {
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
	if scheduled.Spec.Schedule.ConcurrencyPolicy == skprmetav1.ForbidConcurrent && len(active) > 0 {
		return result, nil
	}

	// ...or instruct us to replace existing ones...
	if scheduled.Spec.Schedule.ConcurrencyPolicy == skprmetav1.ReplaceConcurrent {
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
