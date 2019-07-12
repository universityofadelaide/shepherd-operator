package restore

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"time"

	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	errorspkg "github.com/pkg/errors"
	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	v1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-test/deep"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/sync"
	resticutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/restic"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "restore-controller"
)

// Add creates a new Restore Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRestore{
		Config: mgr.GetConfig(),
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
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
	err = c.Watch(&source.Kind{Type: &extensionv1.Restore{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in a Job owned by a Restore
	return c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionv1.Restore{},
	})
}

var _ reconcile.Reconciler = &ReconcileRestore{}

// ReconcileRestore reconciles a Restore object
type ReconcileRestore struct {
	client.Client
	Config *rest.Config
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Restore object and makes changes based on the state read
// and what is in the Restore.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Jobs.
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get
// +kubebuilder:rbac:groups=extension.shepherd,resources=restores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=restores/status,verbs=get;update;patch
func (r *ReconcileRestore) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	restore := &extensionv1.Restore{}

	err := r.Get(context.TODO(), request.NamespacedName, restore)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	backup := &extensionv1.Backup{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      restore.Spec.BackupName,
		Namespace: restore.ObjectMeta.Namespace,
	}, backup)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	switch backup.Status.Phase {
	case v1.PhaseFailed:
		log.Info(fmt.Sprintf("Skipping restore %s because the backup %s failed", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	case v1.PhaseNew:
		// Requeue the operation for 60 seconds if the backup is new.
		return requeueAfterSeconds(60), nil
	case v1.PhaseInProgress:
		// Requeue the operation for 30 seconds if the backup is still in progress.
		return requeueAfterSeconds(30), nil
	}
	// Catch-all for any other non Completed phases.
	if backup.Status.Phase != v1.PhaseCompleted {
		return reconcile.Result{}, nil
	}

	if _, found := restore.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the restore failed
		log.Info(fmt.Sprintf("Restore %s doesn't have a site label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}
	// TODO: Add environment to spec so we don't have to derive the deploymentconfig name.
	if _, found := restore.ObjectMeta.GetLabels()["environment"]; !found {
		log.Info(fmt.Sprintf("Restore %s doesn't have a environment label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	v1client, err := osv1client.NewForConfig(r.Config)
	if err != nil {
		return reconcile.Result{}, errorspkg.Wrap(err, "failed to get deploymentconfig client")
	}
	dcName := fmt.Sprintf("node-%s", restore.ObjectMeta.GetLabels()["environment"])
	dc, err := v1client.DeploymentConfigs(restore.ObjectMeta.Namespace).Get(dcName, metav1.GetOptions{})
	if err != nil {
		// Don't throw an error here to account for restores that were created before an environment was deleted.
		return reconcile.Result{}, nil
	}

	var params = resticutils.PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "docker.io/restic/restic:0.9.5",
		MySQLImage:  "previousnext/mysql",
		WorkingDir:  "/home/shepherd",
		Tags:        []string{},
	}
	spec, err := resticutils.PodSpecRestore(restore, dc, backup.Status.ResticID, params, restore.ObjectMeta.GetLabels()["site"])
	if err != nil {
		return reconcile.Result{}, err
	}

	var (
		parallelism    int32 = 1
		completions    int32 = 1
		activeDeadline int64 = 3600
		backOffLimit   int32 = 2
	)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", resticutils.Prefix, restore.ObjectMeta.Name),
			Namespace: restore.ObjectMeta.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism:           &parallelism,
			Completions:           &completions,
			ActiveDeadlineSeconds: &activeDeadline,
			BackoffLimit:          &backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: spec,
			},
		},
	}

	log.Info("Syncing Job")
	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, job, sync.Job(restore, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Job %s with status: %s", job.ObjectMeta.Name, result)

	log.Info("Syncing status")
	status := extensionv1.RestoreStatus{
		Phase:          v1.PhaseNew,
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	if job.Status.Active > 0 {
		status.Phase = v1.PhaseInProgress
	} else {
		if job.Status.Succeeded > 0 {
			status.Phase = v1.PhaseCompleted
		}
		if job.Status.Failed > 0 {
			status.Phase = v1.PhaseFailed
		}
	}

	if diff := deep.Equal(restore.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		restore.Status = status

		err := r.Status().Update(context.TODO(), restore)
		if err != nil {
			return reconcile.Result{}, errorspkg.Wrap(err, "failed to update status")
		}
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// requeueAfterSeconds returns a reconcile.Result to requeue after seconds time.
func requeueAfterSeconds(seconds int64) reconcile.Result {
	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(seconds) * time.Second,
	}
}
