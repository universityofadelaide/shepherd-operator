package restore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	plog "github.com/prometheus/common/log"
	"github.com/skpr/operator/pkg/utils/k8s/events"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/skpr/operator/internal/controller/backup/restic"
	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	jobutils "github.com/skpr/operator/pkg/utils/k8s/job"
	"github.com/skpr/operator/pkg/utils/k8s/sync"
	resticutils "github.com/skpr/operator/pkg/utils/restic"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "restore-controller"
)

// Add creates a new Restore Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params restic.Params) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params restic.Params) reconcile.Reconciler {
	return &ReconcileRestore{
		Client:   mgr.GetClient(),
		Config:   mgr.GetConfig(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder(ControllerName),
		params:   params,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Restore
	err = c.Watch(&source.Kind{Type: &extensionsv1beta1.Restore{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in a Job owned by a Restore
	return c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionsv1beta1.Restore{},
	})
}

var _ reconcile.Reconciler = &ReconcileRestore{}

// ReconcileRestore reconciles a Restore object
type ReconcileRestore struct {
	client.Client
	Config   *rest.Config
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	params   restic.Params
}

// Reconcile reads that state of the cluster for a Restore object and makes changes based on the state read
// and what is in the Restore.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=restores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=restores/status,verbs=get;update;patch
func (r *ReconcileRestore) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	restore := &extensionsv1beta1.Restore{}

	err := r.Get(context.TODO(), request.NamespacedName, restore)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	backup, err := r.getBackupFromRestore(restore)
	if err != nil {
		message := fmt.Sprintf("Skipping restore %s because backup %s could not be found", restore.ObjectMeta.Name, restore.Spec.BackupName)
		log.Info(message)
		r.recorder.Event(restore, corev1.EventTypeNormal, events.EventCreate, message)
		statusErr := r.syncStatus(restore, extensionsv1beta1.RestoreStatus{Phase: skprmetav1.PhaseFailed}, log)
		if statusErr != nil {
			return reconcile.Result{}, err
		}
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	switch backup.Status.Phase {
	case skprmetav1.PhaseFailed:
		message := fmt.Sprintf("Skipping restore %s because the backup %s failed", restore.ObjectMeta.Name, backup.ObjectMeta.Name)
		log.Info(message)
		r.recorder.Event(restore, corev1.EventTypeNormal, events.EventCreate, message)
		statusErr := r.syncStatus(restore, extensionsv1beta1.RestoreStatus{Phase: skprmetav1.PhaseFailed}, log)
		if statusErr != nil {
			return reconcile.Result{}, statusErr
		}
		return reconcile.Result{}, nil
	case skprmetav1.PhaseUnknown:
		// Requeue the operation for 60 seconds if the backup is new.
		return requeueAfterSeconds(60), nil
	case skprmetav1.PhaseInProgress:
		// Requeue the operation for 30 seconds if the backup is still in progress.
		return requeueAfterSeconds(30), nil
	}
	// Catch-all for any other non Completed phases.
	if backup.Status.Phase != skprmetav1.PhaseCompleted {
		return reconcile.Result{}, nil
	}

	spec, err := resticutils.PodSpecRestore(restore, backup.Status.ResticID, r.params.Pod)
	if err != nil {
		return reconcile.Result{}, err
	}

	job, err := jobutils.NewFromPod(metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", resticutils.RestorePrefix, restore.ObjectMeta.Name),
		Namespace: restore.ObjectMeta.Namespace,
	}, spec)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Syncing Job")

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, job, sync.Job(restore, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Job %s with status: %s", job.ObjectMeta.Name, result)

	log.Info("Syncing status")
	status := extensionsv1beta1.RestoreStatus{
		Phase:          jobutils.GetPhase(job.Status),
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	err = r.syncStatus(restore, status, log)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to update status")
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// syncStatus syncs the status of a restore object.
func (r *ReconcileRestore) syncStatus(restore *extensionsv1beta1.Restore, status extensionsv1beta1.RestoreStatus, log plog.Logger) error {
	if diff := deep.Equal(restore.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		restore.Status = status

		return r.Status().Update(context.TODO(), restore)
	}
	return nil
}

// getBackupFromRestore loads a backup object from a restore.
func (r *ReconcileRestore) getBackupFromRestore(restore *extensionsv1beta1.Restore) (*extensionsv1beta1.Backup, error) {
	backup := &extensionsv1beta1.Backup{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      restore.Spec.BackupName,
		Namespace: restore.ObjectMeta.Namespace,
	}, backup)

	return backup, err
}

// requeueAfterSeconds returns a reconcile.Result to requeue after seconds time.
func requeueAfterSeconds(seconds int64) reconcile.Result {
	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(seconds) * time.Second,
	}
}
