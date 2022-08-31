/*
Copyright 2022.

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

package restore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-test/deep"
	osv1 "github.com/openshift/api/apps/v1"
	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpdmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	metautils "github.com/universityofadelaide/shepherd-operator/internal/k8s/metadata"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "sync-restore-controller"
)

// Reconciler reconciles a Sync object
type Reconciler struct {
	client.Client
	OpenShift osv1client.AppsV1Interface
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Params    Params
}

// Params used by this controller.
type Params struct {
	// Used to filter Backup objects by a key and value pair.
	FilterByLabelAndValue FilterByLabelAndValue
}

// FilterByLabelAndValue is used to filter Backup objects by a key and value pair.
type FilterByLabelAndValue struct {
	Key   string
	Value string
}

//+kubebuilder:rbac:groups=extension.shepherd,resources=restores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=syncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=syncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=syncs/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	sync := &extensionv1.Sync{}

	err := r.Get(ctx, req.NamespacedName, sync)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if !metautils.HasLabelWithValue(sync.ObjectMeta.Labels, r.Params.FilterByLabelAndValue.Key, r.Params.FilterByLabelAndValue.Value) {
		return reconcile.Result{}, nil
	}

	if sync.Status.Backup.Phase != shpdmetav1.PhaseCompleted {
		logger.Info("Skipping. Backup hasn't finished.", "name", sync.Status.Backup.Name, "phase", sync.Status.Backup.Phase)
		return reconcile.Result{}, nil
	}

	dc, err := r.OpenShift.DeploymentConfigs(sync.ObjectMeta.Namespace).Get(ctx, fmt.Sprintf("node-%s", sync.Spec.RestoreEnv), metav1.GetOptions{})
	if err != nil {
		// Don't throw an error here to account for syncs that were created before an environment was deleted.
		return reconcile.Result{}, nil
	}

	// Check the deployment is running so we can create a restore.
	if !isDeploymentRunning(dc) {
		logger.Info("Deployment not yet running, will requeue after 10 seconds")
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, nil
	}

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
			BackupName: sync.Status.Backup.Name,
			Volumes:    sync.Spec.RestoreSpec.Volumes,
			MySQL:      sync.Spec.RestoreSpec.MySQL,
		},
	}

	logger.Info("Syncing Restore")

	if err := controllerutil.SetControllerReference(sync, restore, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.Create(ctx, restore); client.IgnoreNotFound(err) != nil {
		return reconcile.Result{}, err
	}

	if err = r.Get(ctx, types.NamespacedName{
		Namespace: restore.ObjectMeta.Namespace,
		Name:      restore.ObjectMeta.Name,
	}, restore); err != nil {
		return ctrl.Result{}, err
	}

	status := extensionv1.SyncStatusRestore{
		Name:           restore.ObjectMeta.Name,
		Phase:          restore.Status.Phase,
		CompletionTime: restore.Status.StartTime,
	}

	if diff := deep.Equal(sync.Status.Restore, status); diff != nil {
		logger.Info(fmt.Sprintf("Status change dectected: %s", diff))

		sync.Status.Restore = status

		err := r.Status().Update(ctx, sync)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to update status: %w", err)
		}
	}

	logger.Info("Reconcile finished")

	return ctrl.Result{}, nil
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

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Sync{}).
		Owns(&extensionv1.Restore{}).
		Complete(r)
}
