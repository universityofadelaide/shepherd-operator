package shared

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	"github.com/skpr/operator/pkg/mysql"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/k8s/events"
	"github.com/skpr/operator/pkg/utils/random"
	"github.com/skpr/operator/pkg/utils/slice"
)

const (
	// Finalizer used to trigger a deletion of the user prior to the object being deleted.
	Finalizer = "databases.mysql.skpr.io"
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "mysql-database-shared"
)

// Add creates a new Database Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, name string, conn Connection) error {
	return add(mgr, newReconciler(mgr, name, conn))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, name string, conn Connection) reconcile.Reconciler {
	return &ReconcileDatabase{
		ProvisionerName: name,
		Client:          mgr.GetClient(),
		recorder:        mgr.GetRecorder(ControllerName),
		scheme:          mgr.GetScheme(),
		Connection:      conn,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Database
	return c.Watch(&source.Kind{Type: &mysqlv1beta1.Database{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileDatabase{}

// ReconcileDatabase reconciles a Database object
type ReconcileDatabase struct {
	ProvisionerName string
	client.Client
	recorder   record.EventRecorder
	scheme     *runtime.Scheme
	Connection Connection
}

// Connection details users for provisioning databases, users and grants.
type Connection struct {
	Hostname string
	Port     int
	Username string
	Password string
}

// Reconcile reads that state of the cluster for a Database object and makes changes based on the state read
// and what is in the Database.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=databaseclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=databaseclaims/status,verbs=get;update;patch
func (r *ReconcileDatabase) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	database := &mysqlv1beta1.Database{}

	err := r.Get(context.TODO(), request.NamespacedName, database)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// @todo, Can we filter upstream to avoid reconciles?
	if database.Spec.Provisioner != r.ProvisionerName {
		log.Info("Skipping because database is not set to this provisioner:", database.Spec.Provisioner)
		return reconcile.Result{}, nil
	}

	client, err := mysql.New(r.Connection.Hostname, r.Connection.Username, r.Connection.Password, r.Connection.Port)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to setup MySQL client")
	}

	// https://book.kubebuilder.io/beyond_basics/using_finalizers.html
	if database.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !slice.Contains(database.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Adding finalizer:", Finalizer)

			database.ObjectMeta.Finalizers = append(database.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), database); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		// The object is being deleted, ensure that we have the finalizer and delete the database / grant.
		if slice.Contains(database.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Deleting database")

			// our finalizer is present, so lets handle our external dependency
			if err := r.DeleteGrant(client, database, database.Status.Connection); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.DeleteUser(client, database, database.Status.Connection); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.DeleteDatabase(client, database, database.Status.Connection); err != nil {
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			database.ObjectMeta.Finalizers = slice.Remove(database.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), database); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}

		return reconcile.Result{}, nil
	}

	name := hash(fmt.Sprintf("%s_%s", database.ObjectMeta.Namespace, database.ObjectMeta.Name))

	status := mysqlv1beta1.DatabaseStatus{
		ObservedGeneration: database.Generation,
		Connection: mysqlv1beta1.DatabaseStatusConnection{
			Hostname: r.Connection.Hostname,
			Port:     r.Connection.Port,
			Database: name,
			Username: name,
			Password: random.String(16),
		},
	}

	if database.Status.Connection.Hostname != "" {
		status.Connection.Hostname = database.Status.Connection.Hostname
	}

	if database.Status.Connection.Port != 0 {
		status.Connection.Port = database.Status.Connection.Port
	}

	if database.Status.Connection.Database != "" {
		status.Connection.Database = database.Status.Connection.Database
	}

	if database.Status.Connection.Username != "" {
		status.Connection.Username = database.Status.Connection.Username
	}

	if database.Status.Connection.Password != "" {
		status.Connection.Password = database.Status.Connection.Password
	}

	log.Info("Syncing database")

	err = r.SyncDatabase(client, database, status.Connection)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, err
	}

	log.Info("Syncing database user")

	err = r.SyncUser(client, database, status.Connection)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, err
	}

	log.Info("Syncing database grant")

	err = r.SyncGrant(client, database, status.Connection)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, err
	}

	status.Phase = mysqlv1beta1.PhaseReady

	log.Info("Syncing status")

	err = r.SyncStatus(database, status)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, err
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// SyncStatus for the Database object.
func (r *ReconcileDatabase) SyncStatus(database *mysqlv1beta1.Database, status mysqlv1beta1.DatabaseStatus) error {
	if diff := deep.Equal(database.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		database.Status = status

		err := r.Status().Update(context.TODO(), database)
		if err != nil {
			return errors.Wrap(err, "failed to update status")
		}
	}

	return nil
}

// SyncDatabase for the reconcile loop.
func (r *ReconcileDatabase) SyncDatabase(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.Database().Exists(conn.Database)
	if err != nil {
		return errors.Wrap(err, "failed to list")
	}

	if exists {
		return nil
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "Creating database: %s", conn.Database)

	err = client.Database().Create(conn.Database)
	if err != nil {
		return errors.Wrap(err, "failed to create")
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "Database created: %s", conn.Database)

	return nil
}

// DeleteDatabase for the reconcile loop.
func (r *ReconcileDatabase) DeleteDatabase(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.Database().Exists(conn.Database)
	if err != nil {
		return errors.Wrap(err, "failed to list")
	}

	if !exists {
		return nil
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "Deleting database: %s", conn.Database)

	err = client.Database().Delete(conn.Database)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "Database deleted: %s", conn.Database)

	return nil
}

// SyncUser for the reconcile loop.
func (r *ReconcileDatabase) SyncUser(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.User().Exists(conn.Username)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "Creating user: %s", conn.Username)

	err = client.User().Create(conn.Username, conn.Password)
	if err != nil {
		return err
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "User created: %s", conn.Username)

	return nil
}

// DeleteUser for the reconcile loop.
func (r *ReconcileDatabase) DeleteUser(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.User().Exists(conn.Username)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "Deleting user: %s", conn.Username)

	err = client.User().Delete(conn.Username)
	if err != nil {
		return err
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "User deleted: %s", conn.Username)

	return nil
}

// SyncGrant for the reconcile loop.
func (r *ReconcileDatabase) SyncGrant(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.Grant().Exists(conn.Username, conn.Database)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// @todo, Do we need to check if the user exists first?

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "Granting privileges for user: %s", conn.Username)

	err = client.Grant().Create(conn.Username, conn.Database, database.Spec.Privileges)
	if err != nil {
		return err
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventCreate, "User has been granted privileges: %s", conn.Username)

	return nil
}

// DeleteGrant for the reconcile loop.
func (r *ReconcileDatabase) DeleteGrant(client *mysql.Client, database *mysqlv1beta1.Database, conn mysqlv1beta1.DatabaseStatusConnection) error {
	exists, err := client.Grant().Exists(conn.Username, conn.Database)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "Revoking grant for: %s/%s", conn.Username, conn.Database)

	err = client.Grant().Revoke(conn.Username, conn.Database)
	if err != nil {
		return errors.Wrap(err, "failed to revoke grant")
	}

	r.recorder.Eventf(database, corev1.EventTypeNormal, events.EventDelete, "Grant revoked for: %s/%s", conn.Username, conn.Database)

	return nil
}
