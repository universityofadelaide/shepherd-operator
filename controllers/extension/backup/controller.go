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
	"io/ioutil"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	v1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/internal/events"
	jobutils "github.com/universityofadelaide/shepherd-operator/internal/k8s/job"
	"github.com/universityofadelaide/shepherd-operator/internal/restic"
	sliceutils "github.com/universityofadelaide/shepherd-operator/internal/slice"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "backup-restic-controller"
	// Finalizer used by this controller to perform a final operation on object deletion.
	Finalizer = "backups.finalizers.shepherd"
)

// Reconciler reconciles a Backup object
type Reconciler struct {
	client.Client
	Config   *rest.Config
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	Params   Params
}

// Params which are provided to this controller.
type Params struct {
	// Parameters which are used when provisioning a Pod instance.
	PodSpec restic.PodSpecParams
}

//+kubebuilder:rbac:groups=v1,resources=pods,verbs=get;list
//+kubebuilder:rbac:groups=v1,resources=pods/log,verbs=get
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconcile loop")

	backup := &extensionv1.Backup{}

	err := r.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if backup.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, if it does not have this finalizer,
		// add it and update the object.
		if !sliceutils.Contains(backup.ObjectMeta.Finalizers, Finalizer) {
			logger.Info("Adding finalizer")

			backup.ObjectMeta.Finalizers = append(backup.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(ctx, backup); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		// The object is being deleted, ensure that the finalizer exists then
		// create a job to delete the restic snapshot.
		if sliceutils.Contains(backup.ObjectMeta.Finalizers, Finalizer) {
			// Check status of restic-delete job.
			newJob := &batchv1.Job{}

			nsn := types.NamespacedName{
				Namespace: backup.Namespace,
				Name:      fmt.Sprintf("%s-delete-%s", restic.Prefix, backup.Name),
			}

			err := r.Get(ctx, nsn, newJob)
			if err != nil {
				if !kerrors.IsNotFound(err) {
					return reconcile.Result{Requeue: true}, err
				}

				logger.Info("Forgetting the Restic snapshot", "id", backup.Status.ResticID)

				if backup.Status.ResticID == "" {
					// Allow the backup to delete when we don't know the restic id.
					logger.Info("No restic ID associated when attempting to delete backup", "name", backup.ObjectMeta.Name)
					return r.removeFinalizer(ctx, backup)
				} else {
					// Job doesnt exist, create it.
					err := r.DeleteResticSnapshot(ctx, backup)
					return reconcile.Result{RequeueAfter: 5 * time.Second}, err
				}
			}

			if jobutils.IsFinished(newJob) {
				logger.Info("Removing finalizer", "finalizer", Finalizer)
				return r.removeFinalizer(ctx, backup)
			}

			logger.Info("Requeuing to wait for finalizer Job to finish")

			return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
		}

		return reconcile.Result{}, nil
	}

	// Backup has completed or failed, return early.
	if backup.Status.Phase == v1.PhaseCompleted || backup.Status.Phase == v1.PhaseFailed {
		return reconcile.Result{}, nil
	}

	if _, found := backup.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the backup failed
		logger.Info(fmt.Sprintf("Backup %s doesn't have a site label, skipping.", backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	return r.SyncJob(ctx, logger, backup)
}

// SyncJob creates or updates the restic backup jobs.
func (r *Reconciler) SyncJob(ctx context.Context, log logr.Logger, backup *extensionv1.Backup) (reconcile.Result, error) {
	// Backup has completed or failed, return early.
	if backup.Status.Phase == v1.PhaseCompleted || backup.Status.Phase == v1.PhaseFailed {
		return reconcile.Result{}, nil
	}

	spec, err := restic.PodSpecBackup(backup, r.Params.PodSpec, backup.ObjectMeta.GetLabels()["site"])
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
			Name:      fmt.Sprintf("%s-%s", restic.Prefix, backup.ObjectMeta.Name),
			Namespace: backup.ObjectMeta.Namespace,
			Labels: map[string]string{
				"app":          "restic",
				"resticAction": "backup",
			},
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
			resticId, err := getResticIdFromJob(ctx, r.Config, job)
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

		err := r.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	log.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// DeleteResticSnapshot creates the job to forget a restic snapshot.
func (r *Reconciler) DeleteResticSnapshot(ctx context.Context, backup *extensionv1.Backup) error {
	if backup.Status.ResticID == "" {
		return errors.Errorf("Could't delete restic snapshot. Restic ID missing for backup: %s", backup.ObjectMeta.Name)
	}

	spec, err := restic.PodSpecDelete(
		backup.Status.ResticID,
		backup.ObjectMeta.Namespace,
		backup.ObjectMeta.GetLabels()["site"],
		r.Params.PodSpec,
	)
	if err != nil {
		return err
	}

	var (
		parallelism    int32 = 1
		completions    int32 = 1
		activeDeadline int64 = 3600
		backOffLimit   int32 = 2
		//ttl            int32 = 3600
	)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-delete-%s", restic.Prefix, backup.Name),
			Namespace: backup.ObjectMeta.Namespace,
			Labels: map[string]string{
				"app":          "restic",
				"resticAction": "delete",
			},
		},
		Spec: batchv1.JobSpec{
			// @todo uncomment this when the feature becomes available (requires kube v1.12+).
			// ttlSecondsAfterFinished: &ttl,
			Parallelism:           &parallelism,
			Completions:           &completions,
			ActiveDeadlineSeconds: &activeDeadline,
			BackoffLimit:          &backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: spec,
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, job, func() error {
		return nil
	})
	if err != nil {
		return err
	}

	switch result {
	case controllerutil.OperationResultCreated:
		r.Recorder.Eventf(backup, corev1.EventTypeNormal, events.EventCreate, "Job has been created: %s", job.ObjectMeta.Name)
	case controllerutil.OperationResultUpdated:
		r.Recorder.Eventf(backup, corev1.EventTypeNormal, events.EventUpdate, "Job has been updated: %s", job.ObjectMeta.Name)
	}

	return err
}

// getResticIdFromJob parses output from a job's pods and returns a restic ID from the logs.
func getResticIdFromJob(ctx context.Context, config *rest.Config, job *batchv1.Job) (string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	pods, err := clientset.CoreV1().Pods(job.ObjectMeta.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(job.Spec.Template.ObjectMeta.Labels).String(),
	})
	if err != nil {
		return "", err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodSucceeded {
			continue
		}

		podLogs, err := getPodLogs(ctx, clientset, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
		if err != nil {
			return "", err
		}

		resticId := restic.ParseSnapshotID(podLogs)
		if resticId != "" {
			return resticId, nil
		}
	}

	return "", nil
}

// getPodLogs gets the logs from the restic container from a pod as a string.
func getPodLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace string, podName string) (string, error) {
	body, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: restic.ResticBackupContainerName,
	}).Stream(ctx)
	if err != nil {
		return "", err
	}
	defer body.Close()

	podLogs, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(podLogs), nil
}

// Helper function to remove the finalizer and exit a reconcile loop.
func (r *Reconciler) removeFinalizer(ctx context.Context, backup *extensionv1.Backup) (reconcile.Result, error) {
	backup.ObjectMeta.Finalizers = sliceutils.Remove(backup.ObjectMeta.Finalizers, Finalizer)
	err := r.Update(ctx, backup)
	return reconcile.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Backup{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
