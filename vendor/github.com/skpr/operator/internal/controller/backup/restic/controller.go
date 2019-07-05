package restic

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	k8ssync "github.com/skpr/operator/pkg/utils/k8s/sync"
	"github.com/skpr/operator/pkg/utils/random"
	resticutils "github.com/skpr/operator/pkg/utils/restic"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "backup-restic-controller"

// Add creates a new Backup Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params Params) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params Params) reconcile.Reconciler {
	return &ReconcileBackup{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		params: params,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1beta1.CronJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionsv1beta1.Backup{},
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &extensionsv1beta1.Backup{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileBackup{}

// ReconcileBackup reconciles a Backup object
type ReconcileBackup struct {
	client.Client
	scheme *runtime.Scheme
	params Params
}

// Params which are passed into this controller.
type Params struct {
	Pod     resticutils.PodSpecParams
	CronJob ParamsCronJob
}

// ParamsCronJob used for grouping CronJob configuration.
type ParamsCronJob struct {
	StartingDeadline int64
	ActiveDeadline   int64
	BackoffLimit     int32
	SuccessHistory   int32
	FailedHistory    int32
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=backups/status,verbs=get;update;patch
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	backup := &extensionsv1beta1.Backup{}

	err := r.Get(context.TODO(), request.NamespacedName, backup)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", resticutils.Prefix, backup.ObjectMeta.Name),
			Namespace: backup.ObjectMeta.Namespace,
		},
		Data: map[string][]byte{
			resticutils.ResticPassword: []byte(random.String(32)),
		},
	}

	log.Info("Syncing Secret")

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, secret, k8ssync.Secret(backup, secret.Data, false, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Secret %s with status: %s", secret.ObjectMeta.Name, result)

	spec, err := resticutils.PodSpec(backup, r.params.Pod)
	if err != nil {
		return reconcile.Result{}, err
	}

	var (
		parallelism int32 = 1
		completions int32 = 1
	)

	cronjob := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", resticutils.Prefix, backup.ObjectMeta.Name),
			Namespace: backup.ObjectMeta.Namespace,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:                   backup.Spec.Schedule,
			ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
			StartingDeadlineSeconds:    &r.params.CronJob.StartingDeadline,
			SuccessfulJobsHistoryLimit: &r.params.CronJob.SuccessHistory,
			FailedJobsHistoryLimit:     &r.params.CronJob.FailedHistory,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Parallelism:           &parallelism,
					Completions:           &completions,
					ActiveDeadlineSeconds: &r.params.CronJob.ActiveDeadline,
					BackoffLimit:          &r.params.CronJob.BackoffLimit,
					Template: corev1.PodTemplateSpec{
						Spec: spec,
					},
				},
			},
		},
	}

	log.Info("Syncing CronJob")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, cronjob, k8ssync.CronJob(backup, cronjob.Spec, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced CronJob %s with status: %s", cronjob.ObjectMeta.Name, result)

	status := extensionsv1beta1.BackupStatus{
		LastScheduleTime: cronjob.Status.LastScheduleTime,
	}

	log.Info("Syncing status")

	if diff := deep.Equal(backup.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		backup.Status = status

		err := r.Status().Update(context.TODO(), backup)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}
