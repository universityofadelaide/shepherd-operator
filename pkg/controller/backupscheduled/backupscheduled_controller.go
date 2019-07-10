package backupscheduled

import (
	"context"
	"fmt"
	"time"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	v1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
	"io/ioutil"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/sync"
	resticutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/restic"
	"github.com/universityofadelaide/shepherd-operator/pkg/apis/extension"
	"github.com/gorhill/cronexpr"
)

const ControllerName = "backup-scheduled-controller"

// Add creates a new BackupScheduled Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackupScheduled{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to BackupScheduled
	err = c.Watch(&source.Kind{Type: &extensionv1.BackupScheduled{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Backup
	err = c.Watch(&source.Kind{Type: &extensionv1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileBackupScheduled{}

// ReconcileBackupScheduled reconciles a BackupScheduled object
type ReconcileBackupScheduled struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a BackupScheduled object and makes changes based on the state read
// and what is in the BackupScheduled.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backupscheduleds/status,verbs=get;update;patch
func (r *ReconcileBackupScheduled) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)
	log.Info("Starting reconcile loop")

	backupScheduled := &extensionv1.BackupScheduled{}
	err := r.Get(context.TODO(), request.NamespacedName, backupScheduled)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if _, found := backupScheduled.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the backup failed
		log.Info(fmt.Sprintf("BackupScheduled %s doesn't have a site label, skipping.", backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	now := time.Now()

	// Set up the first schedule if
	if backupScheduled.Status.LastExecutedTime == nil {
		next := cronexpr.MustParse(backupScheduled.Spec.Schedule).Next(now)
		nextDuration := next.Unix() - now.Unix()
		return reconcile.Result{
			RequeueAfter: time.Duration(nextDuration),
		}, nil
	}




	var backup extensionv1.Backup
	backup.Name = fmt.Sprintf("%s-%d", backupScheduled.Name, now.Unix())
	backup.Labels = backupScheduled.Labels
	backup.Spec.MySQL = backupScheduled.Spec.MySQL
	backup.Spec.Volumes = backupScheduled.Spec.Volumes


	var params = resticutils.PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "docker.io/restic/restic:0.9.5",
		MySQLImage:  "docker.io/library/mariadb:10",
		WorkingDir:  "/home/shepherd",
		Tags:        []string{},
	}
	spec, err := resticutils.PodSpecBackup(backup, params, backup.ObjectMeta.GetLabels()["site"])
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
			Name:      fmt.Sprintf("%s-%s", resticutils.Prefix, backup.ObjectMeta.Name),
			Namespace: backup.ObjectMeta.Namespace,
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
	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, job, sync.Job(backup, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Job %s with status: %s", job.ObjectMeta.Name, result)

	log.Info("Syncing status")
	status := extensionv1.BackupStatus{
		Phase:          v1.PhaseNew,
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	if job.Status.Active > 0 {
		status.Phase = v1.PhaseInProgress
	} else {
		if job.Status.Succeeded > 0 {
			resticId, err := getResticIdFromJob(r.Config, job)
			if err != nil {
				return reconcile.Result{}, errors.Wrap(err, "failed to parse resticId")
			}

			if resticId != "" {
				status.ResticID = resticId
				status.Phase = v1.PhaseCompleted
			} else {
				status.Phase = v1.PhaseFailed
			}
		}
		if job.Status.Failed > 0 {
			status.Phase = v1.PhaseFailed
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
