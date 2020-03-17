package imagescheduled

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

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	"github.com/skpr/operator/pkg/utils/clock"
	errutils "github.com/skpr/operator/pkg/utils/controller/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/k8s/events"
	scheduledutils "github.com/skpr/operator/pkg/utils/scheduled"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "mysql-imagescheduleds"
	// OwnerKey used to query for child Images.
	OwnerKey = ".metadata.controller"
)

// Add creates a new ImageScheduled Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileImageScheduled{
		Clock:    clock.New(),
		Client:   mgr.GetClient(),
		recorder: mgr.GetRecorder(ControllerName),
		scheme:   mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = mgr.GetFieldIndexer().IndexField(&mysqlv1beta1.Image{}, OwnerKey, func(rawObj runtime.Object) []string {
		image := rawObj.(*mysqlv1beta1.Image)

		owner := metav1.GetControllerOf(image)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != mysqlv1beta1.SchemeGroupVersion.String() || owner.Kind != "ImageScheduled" {
			return nil
		}

		return []string{owner.Name}
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &mysqlv1beta1.ImageScheduled{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileImageScheduled{}

// ReconcileImageScheduled reconciles a ImageScheduled object
type ReconcileImageScheduled struct {
	clock.Clock
	client.Client
	recorder record.EventRecorder
	scheme   *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ImageScheduled object and makes changes based on the state read
// and what is in the ImageScheduled.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=imagescheduleds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=imagescheduleds/status,verbs=get;update;patch
func (r *ReconcileImageScheduled) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	scheduled := &mysqlv1beta1.ImageScheduled{}

	err := r.Get(context.TODO(), request.NamespacedName, scheduled)
	if err != nil {
		return reconcile.Result{}, errutils.IgnoreNotFound(err)
	}

	log.Info("Querying Images")

	var images mysqlv1beta1.ImageList

	err = r.List(context.Background(), client.MatchingField(OwnerKey, request.Name), &images)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Filtering Images")

	active, successful, failed, err := r.SortImages(scheduled, images)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Infof("Job count: active=%d success=%d failed=%d", len(active), len(successful), len(failed))

	err = r.Status().Update(context.TODO(), scheduled)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Cleaning up old Images")

	err = r.Cleanup(scheduled, successful, failed)
	if err != nil {
		return reconcile.Result{}, err
	}

	if scheduled.Spec.Schedule.Suspend != nil && *scheduled.Spec.Schedule.Suspend {
		log.Info("Scheduling has been suspended. Skipping.")
		return reconcile.Result{}, nil
	}

	return r.ScheduleNextImage(scheduled, active)
}

// SortImages into active, successful, failed.
func (r *ReconcileImageScheduled) SortImages(scheduled *mysqlv1beta1.ImageScheduled, images mysqlv1beta1.ImageList) ([]*mysqlv1beta1.Image, []*mysqlv1beta1.Image, []*mysqlv1beta1.Image, error) {
	var (
		active       []*mysqlv1beta1.Image
		successful   []*mysqlv1beta1.Image
		failed       []*mysqlv1beta1.Image
		lastSchedule *time.Time
	)

	for i, image := range images.Items {
		switch image.Status.Phase {
		case skprmetav1.PhaseFailed:
			failed = append(failed, &images.Items[i])
		case skprmetav1.PhaseCompleted:
			successful = append(successful, &images.Items[i])
		default:
			// If it doesn't have a status assigned yet then the Image is most likely still be created.
			active = append(active, &images.Items[i])
		}

		// We'll store the launch time in an annotation, so we'll reconstitute that from
		// the active Images themselves.
		scheduledTime, err := scheduledutils.GetScheduledTime(image.Annotations)
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

// Cleanup old successful and failed Images.
func (r *ReconcileImageScheduled) Cleanup(scheduled *mysqlv1beta1.ImageScheduled, successful, failed []*mysqlv1beta1.Image) error {
	if scheduled.Spec.Schedule.FailedHistoryLimit != nil {
		sort.Slice(failed, func(i, j int) bool {
			if failed[i].Status.StartTime == nil {
				return failed[j].Status.StartTime != nil
			}

			return failed[i].Status.StartTime.Before(failed[j].Status.StartTime)
		})
		for i, image := range failed {
			if int32(i) >= int32(len(failed))-*scheduled.Spec.Schedule.FailedHistoryLimit {
				break
			}

			if err := r.Delete(context.Background(), image, client.PropagationPolicy(metav1.DeletePropagationBackground)); errutils.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old failed job")
			}

			r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting failed Image: %s", image.ObjectMeta.Name)
		}
	}

	if scheduled.Spec.Schedule.SuccessfulHistoryLimit != nil {
		sort.Slice(successful, func(i, j int) bool {
			if successful[i].Status.StartTime == nil {
				return successful[j].Status.StartTime != nil
			}

			return successful[i].Status.StartTime.Before(successful[j].Status.StartTime)
		})

		for i, image := range successful {
			if int32(i) >= int32(len(successful))-*scheduled.Spec.Schedule.SuccessfulHistoryLimit {
				break
			}

			if err := r.Delete(context.Background(), image, client.PropagationPolicy(metav1.DeletePropagationBackground)); errutils.IgnoreNotFound(err) != nil {
				return errors.Wrap(err, "failed to delete old successful job")
			}

			r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventDelete, "Deleting successful Image: %s", image.ObjectMeta.Name)
		}
	}

	return nil
}

// ScheduleNextImage checks if a new Image should be created.
func (r *ReconcileImageScheduled) ScheduleNextImage(scheduled *mysqlv1beta1.ImageScheduled, active []*mysqlv1beta1.Image) (reconcile.Result, error) {
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

	image, err := buildImage(scheduled, r.scheme, missedRun)
	if err != nil {
		return result, errors.Wrap(err, "failed to generate Image")
	}

	r.recorder.Eventf(scheduled, corev1.EventTypeNormal, events.EventCreate, "Creating image: %s", image.ObjectMeta.Name)

	if err := r.Create(context.Background(), image); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create Image")
	}

	return result, nil
}
