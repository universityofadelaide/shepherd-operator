/*
Copyright 2019 University of Adelaide.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sync

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	osv1 "github.com/openshift/api/apps/v1"
	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	errorspkg "github.com/pkg/errors"
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

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/controller/logger"
	syncutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/sync"
	resticutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/restic"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "sync-controller"

// Add creates a new Sync Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	v1client, err := osv1client.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	return add(mgr, newReconciler(mgr, v1client))
}

// newReconciler returns a new ReconcileSync.
func newReconciler(mgr manager.Manager, osclient osv1client.AppsV1Interface) reconcile.Reconciler {
	return &ReconcileSync{
		OsClient: osclient,
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sync-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Sync
	err = c.Watch(&source.Kind{Type: &extensionv1.Sync{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in a Backup owned by a Sync.
	err = c.Watch(&source.Kind{Type: &extensionv1.Backup{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionv1.Sync{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in a Restore owned by a Sync.
	return c.Watch(&source.Kind{Type: &extensionv1.Restore{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionv1.Sync{},
	})
}

var _ reconcile.Reconciler = &ReconcileSync{}

// ReconcileSync reconciles a Sync object
type ReconcileSync struct {
	client.Client
	OsClient osv1client.AppsV1Interface
	scheme   *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Sync object and makes changes based on the state read
// and what is in the Sync.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=extension.shepherd,resources=syncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=syncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get
// +kubebuilder:rbac:groups=extension.shepherd,resources=restores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=restores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=restores/finalizers,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileSync) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)
	log.Info("Starting reconcile loop")

	sync := &extensionv1.Sync{}
	err := r.Get(context.TODO(), request.NamespacedName, sync)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	backup := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("sync-%s-backup", sync.ObjectMeta.Name),
			Namespace: sync.ObjectMeta.Namespace,
			Labels: map[string]string{
				"site":        sync.Spec.Site,
				"environment": sync.Spec.BackupEnv,
			},
		},
		Spec: sync.Spec.BackupSpec,
	}

	log.Info("Syncing Backup")
	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, backup, syncutils.Backup(sync, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Backup %s with status: %s", backup.ObjectMeta.Name, result)

	status := extensionv1.SyncStatus{
		BackupName:     backup.ObjectMeta.Name,
		BackupPhase:    backup.Status.Phase,
		StartTime:      backup.Status.StartTime,
		CompletionTime: backup.Status.CompletionTime,
	}

	dc, err := r.OsClient.DeploymentConfigs(sync.ObjectMeta.Namespace).Get(fmt.Sprintf("node-%s", sync.Spec.RestoreEnv), metav1.GetOptions{})
	if err != nil {
		// Don't throw an error here to account for syncs that were created before an environment was deleted.
		return reconcile.Result{}, nil
	}

	// Check the deployment is running so we can create a restore.
	ret := reconcile.Result{}
	if !isDeploymentRunning(dc) {
		log.Info("Deployment not yet running, will requeue after 10 seconds")
		ret = resticutils.RequeueAfterSeconds(10)
	} else {
		restore := &extensionv1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("sync-%s-restore", sync.ObjectMeta.Name),
				Namespace: sync.ObjectMeta.Namespace,
				Labels: map[string]string{
					"site":        sync.Spec.Site,
					"environment": sync.Spec.RestoreEnv,
				},
			},
			Spec: extensionv1.RestoreSpec{
				BackupName: backup.ObjectMeta.Name,
				Volumes:    sync.Spec.RestoreSpec.Volumes,
				MySQL:      sync.Spec.RestoreSpec.MySQL,
			},
		}

		log.Info("Syncing Restore")
		result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, restore, syncutils.Restore(sync, r.scheme))
		if err != nil {
			return reconcile.Result{}, err
		}
		log.Infof("Synced Restore %s with status: %s", restore.ObjectMeta.Name, result)

		status.RestoreName = restore.ObjectMeta.Name
		status.RestorePhase = restore.Status.Phase
		// Use the restore completion time if it's after the backup completion time.
		if status.CompletionTime == nil || (restore.Status.CompletionTime != nil && restore.Status.CompletionTime.After(status.CompletionTime.Time)) {
			status.CompletionTime = restore.Status.CompletionTime
		}
	}

	if diff := deep.Equal(sync.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		sync.Status = status

		err := r.Status().Update(context.TODO(), sync)
		if err != nil {
			return reconcile.Result{}, errorspkg.Wrap(err, "failed to update status")
		}
	}

	log.Info("Reconcile finished")
	return ret, nil
}

// isDeploymentRunning checks if a deployment config is available and running.
func isDeploymentRunning(dc *osv1.DeploymentConfig) bool {
	for _, condition := range dc.Status.Conditions {
		if condition.Type == osv1.DeploymentAvailable {
			return condition.Status == corev1.ConditionTrue
		}
	}

	return false
}
