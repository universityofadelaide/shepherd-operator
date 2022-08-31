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

package backup

import (
	"context"
	"fmt"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
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
	ControllerName = "sync-backup-controller"
)

// Reconciler reconciles a Sync object
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Params   Params
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

//+kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/finalizers,verbs=get;list;watch;create;update;patch;delete
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

	if sync.Status.Backup.Phase == shpdmetav1.PhaseCompleted || sync.Status.Backup.Phase == shpdmetav1.PhaseFailed {
		logger.Info("Skipping. Backup has finished.", "name", sync.Status.Backup.Name, "phase", sync.Status.Backup.Phase)
		return reconcile.Result{}, nil
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

	if err := controllerutil.SetControllerReference(sync, backup, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Syncing Backup")

	if err := r.Create(ctx, backup); client.IgnoreNotFound(err) != nil {
		return reconcile.Result{}, err
	}

	if err = r.Get(ctx, types.NamespacedName{
		Namespace: backup.ObjectMeta.Namespace,
		Name:      backup.ObjectMeta.Name,
	}, backup); err != nil {
		return ctrl.Result{}, err
	}

	status := extensionv1.SyncStatusBackup{
		Name:      backup.ObjectMeta.Name,
		Phase:     backup.Status.Phase,
		StartTime: backup.Status.StartTime,
	}

	if diff := deep.Equal(sync.Status.Backup, status); diff != nil {
		logger.Info(fmt.Sprintf("Status change dectected: %s", diff))

		sync.Status.Backup = status

		err := r.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	logger.Info("Reconcile finished")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Sync{}).
		Owns(&extensionv1.Backup{}).
		Complete(r)
}
