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

	"github.com/go-test/deep"
	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpdmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/internal/events"
	"github.com/universityofadelaide/shepherd-operator/internal/restic"
	resticutils "github.com/universityofadelaide/shepherd-operator/internal/restic"
)

// Reconciler reconciles a Restore object
type Reconciler struct {
	client.Client
	OpenShift osv1client.AppsV1Interface
	Config    *rest.Config
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Params    Params
}

// Params which are provided to this controller.
type Params struct {
	// Parameters which are used when provisioning a Pod instance.
	PodSpec restic.PodSpecParams
}

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconcile loop")

	restore := &extensionv1.Restore{}

	err := r.Get(ctx, req.NamespacedName, restore)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	backup := &extensionv1.Backup{}

	err = r.Get(ctx, types.NamespacedName{
		Name:      restore.Spec.BackupName,
		Namespace: restore.ObjectMeta.Namespace,
	}, backup)
	if err != nil {
		if kerrors.IsNotFound(err) {
			r.Recorder.Eventf(restore, corev1.EventTypeNormal, events.EventError, "Backup not found: %s", restore.Spec.BackupName)
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	switch backup.Status.Phase {
	case shpdmetav1.PhaseFailed:
		logger.Info(fmt.Sprintf("Skipping restore %s because the backup %s failed", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	case shpdmetav1.PhaseNew:
		// Requeue the operation for 30 seconds if the backup is new.
		logger.Info(fmt.Sprintf("Requeueing restore %s because the backup %s is New", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return resticutils.RequeueAfterSeconds(30), nil
	case shpdmetav1.PhaseInProgress:
		// Requeue the operation for 15 seconds if the backup is still in progress.
		logger.Info(fmt.Sprintf("Requeueing restore %s because the backup %s is In Progress", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return resticutils.RequeueAfterSeconds(15), nil
	}

	// Catch-all for any other non Completed phases.
	if backup.Status.Phase != shpdmetav1.PhaseCompleted {
		logger.Info(fmt.Sprintf("Skipping restore %s because the backup %s is in an unknown state: %s", restore.ObjectMeta.Name, backup.ObjectMeta.Name, backup.Status.Phase))
		return reconcile.Result{}, nil
	}

	if _, found := restore.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the restore failed
		logger.Info(fmt.Sprintf("Restore %s doesn't have a site label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}
	// TODO: Add environment to spec so we don't have to derive the deploymentconfig name.
	if _, found := restore.ObjectMeta.GetLabels()["environment"]; !found {
		logger.Info(fmt.Sprintf("Restore %s doesn't have a environment label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	dcName := fmt.Sprintf("node-%s", restore.ObjectMeta.GetLabels()["environment"])

	dc, err := r.OpenShift.DeploymentConfigs(restore.ObjectMeta.Namespace).Get(ctx, dcName, metav1.GetOptions{})
	if err != nil {
		// Don't throw an error here to account for restores that were ted before an environment was deleted.
		return reconcile.Result{}, nil
	}

	spec, err := resticutils.PodSpecRestore(restore, dc, backup.Status.ResticID, r.Params.PodSpec, restore.ObjectMeta.GetLabels()["site"])
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

	logger.Info("Creating Job")

	if err := controllerutil.SetControllerReference(backup, job, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.Create(ctx, job); client.IgnoreNotFound(err) != nil {
		return reconcile.Result{}, err
	}

	if err = r.Get(ctx, types.NamespacedName{
		Namespace: job.ObjectMeta.Namespace,
		Name:      job.ObjectMeta.Name,
	}, job); err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Syncing status")

	status := extensionv1.RestoreStatus{
		Phase:          shpdmetav1.PhaseNew,
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	if job.Status.Active > 0 {
		status.Phase = shpdmetav1.PhaseInProgress
	} else {
		if job.Status.Succeeded > 0 {
			status.Phase = shpdmetav1.PhaseCompleted
		} else if job.Status.Failed > 0 {
			status.Phase = shpdmetav1.PhaseFailed
		}
	}

	if diff := deep.Equal(restore.Status, status); diff != nil {
		logger.Info(fmt.Sprintf("Status change dectected: %s", diff))

		restore.Status = status

		err := r.Status().Update(context.TODO(), restore)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to update status: %w", err)
		}
	}

	logger.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Restore{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
