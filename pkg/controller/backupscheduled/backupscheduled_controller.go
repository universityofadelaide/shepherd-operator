package backupscheduled

import (
	"context"
	"fmt"
	"time"

	"github.com/go-test/deep"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/sync"
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

	if backupScheduled.Spec.Schedule == "" {
		err := errors.New("BackupScheduled doesn't have a schedule.")
		log.Error(err.Error())
		return reconcile.Result{}, err
	}

	if _, found := backupScheduled.ObjectMeta.GetLabels()["site"]; !found {
		err := errors.New("BackupScheduled doesn't have a site label.")
		log.Error(err.Error())
		return reconcile.Result{}, err
	}

	log.Info("Calculating next scheduled backup.")
	now := time.Now()
	status := extensionv1.BackupScheduledStatus{
		LastExecutedTime: &metav1.Time{getScheduleComparison(backupScheduled.Status, now)},
	}

	// Check if we are currently due for a backup to be scheduled.
	next := cronexpr.MustParse(backupScheduled.Spec.Schedule).Next(status.LastExecutedTime.Time)
	if next.Before(now) || next.Equal(now) {
		// Due for a backup - create object.
		log.Info("Backup due - creating.")
		var backup *extensionv1.Backup
		backup.Name = fmt.Sprintf("%s-%d", backupScheduled.Name, now.Unix())
		backup.Labels = backupScheduled.Labels
		backup.Spec.MySQL = backupScheduled.Spec.MySQL
		backup.Spec.Volumes = backupScheduled.Spec.Volumes

		result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, backup, sync.Backup(backup, r.scheme))
		if err != nil {
			return reconcile.Result{
				Requeue: true,
			}, err
		}
		log.Info(fmt.Sprintf("Created backup with result: %s.", result))

		// Update status with LastExecutedTime.
		status.LastExecutedTime.Time = time.Now()
	}

	if diff := deep.Equal(backupScheduled.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		backupScheduled.Status = status

		err := r.Status().Update(context.TODO(), backupScheduled)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	next = cronexpr.MustParse(backupScheduled.Spec.Schedule).Next(now)
	nextDuration := next.Unix() - now.Unix()

	log.Info(fmt.Sprintf("Reconcile finished, requeued for %s", next.Format(time.RFC3339)))
	return reconcile.Result{
		RequeueAfter: time.Duration(nextDuration),
	}, nil
}

func getScheduleComparison(s extensionv1.BackupScheduledStatus, now time.Time) time.Time {
	if s.LastExecutedTime != nil {
		return s.LastExecutedTime.Time
	}

	return now
}