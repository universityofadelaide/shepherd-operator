package restic

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	jobutils "github.com/skpr/operator/pkg/utils/k8s/job"
	"github.com/skpr/operator/pkg/utils/k8s/pod/logs"
	"github.com/skpr/operator/pkg/utils/k8s/sync"
	resticutils "github.com/skpr/operator/pkg/utils/restic"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "backup-restic-controller"

// Add creates a new Backup Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params Params) error {
	r, err := newReconciler(mgr, params)
	if err != nil {
		return err
	}
	return add(mgr, r)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params Params) (reconcile.Reconciler, error) {
	logsClient, err := logs.New(mgr.GetConfig())
	if err != nil {
		return &ReconcileBackup{}, err
	}
	return &ReconcileBackup{
		Client:    mgr.GetClient(),
		LogClient: logsClient,
		scheme:    mgr.GetScheme(),
		params:    params,
	}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
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
	LogClient logs.Interface
	scheme    *runtime.Scheme
	params    Params
}

// Params which are passed into this controller.
type Params struct {
	Pod resticutils.PodSpecParams
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Jobs.
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
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

	spec, err := resticutils.PodSpecBackup(backup, r.params.Pod)
	if err != nil {
		return reconcile.Result{}, err
	}

	job, err := jobutils.NewFromPod(metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", resticutils.Prefix, backup.ObjectMeta.Name),
		Namespace: backup.ObjectMeta.Namespace,
	}, spec)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Syncing Job")

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, job, sync.Job(backup, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Job %s with status: %s", job.ObjectMeta.Name, result)

	log.Info("Syncing Backup status")
	status := extensionsv1beta1.BackupStatus{
		Phase:          jobutils.GetPhase(job.Status),
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	if status.Phase == skprmetav1.PhaseCompleted {
		resticID, err := r.getResticIDFromJob(job)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to parse resticID")
		}
		if resticID != "" {
			status.ResticID = resticID
		} else {
			status.Phase = skprmetav1.PhaseFailed
		}
	}

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

// getResticIDFromJob parses output from a job's pods and returns a restic ID from the logs.
func (r *ReconcileBackup) getResticIDFromJob(job *batchv1.Job) (string, error) {
	var resticID string

	pods := corev1.PodList{}
	err := r.List(context.TODO(), &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(job.Spec.Template.ObjectMeta.Labels)),
		Namespace:     job.ObjectMeta.Namespace,
	}, &pods)
	if err != nil {
		return resticID, err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodSucceeded {
			continue
		}
		podLogs, err := r.LogClient.Get(pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, resticutils.ResticBackupContainerName)
		if err != nil {
			return resticID, err
		}
		resticID = resticutils.ParseSnapshotID(podLogs)
		if resticID != "" {
			return resticID, nil
		}
	}

	return resticID, nil
}
