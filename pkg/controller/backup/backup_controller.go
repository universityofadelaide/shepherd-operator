package backup

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	v1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/controller/logger"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/sync"
	resticutils "github.com/universityofadelaide/shepherd-operator/pkg/utils/restic"
	"github.com/universityofadelaide/shepherd-operator/pkg/utils/slice"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "backup-restic-controller"
const Finalizer = "backups.finalizers.shepherd"

// Add creates a new Backup Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackup{
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

	// Watch for changes to Backup
	err = c.Watch(&source.Kind{Type: &extensionv1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in a Job owned by a Backup.
	return c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &extensionv1.Backup{},
	})
}

var _ reconcile.Reconciler = &ReconcileBackup{}

// ReconcileBackup reconciles a Backup object
type ReconcileBackup struct {
	client.Client
	Config *rest.Config
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Jobs.
// +kubebuilder:rbac:groups=v1,resources=pods,verbs=get;list
// +kubebuilder:rbac:groups=v1,resources=pods/log,verbs=get
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/finalizers,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)
	log.Info("Starting reconcile loop")

	backup := &extensionv1.Backup{}
	err := r.Get(context.TODO(), request.NamespacedName, backup)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if backup.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, if it does not have this finalizer,
		// add it and update the object.
		if !slice.Contains(backup.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Adding finalizer")

			backup.ObjectMeta.Finalizers = append(backup.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), backup); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		// The object is being deleted, ensure that the finalizer exists then
		// create a job to delete the restic snapshot.
		if slice.Contains(backup.ObjectMeta.Finalizers, Finalizer) {
			// Check status of restic-delete job.
			newJob := &batchv1.Job{}
			nsn := types.NamespacedName{
				Namespace: backup.Namespace,
				Name:      fmt.Sprintf("%s-delete-%s", resticutils.Prefix, backup.Name),
			}
			err := r.Get(context.TODO(), nsn, newJob)
			if err != nil {
				if !kerrors.IsNotFound(err) {
					return reconcile.Result{Requeue: true}, err
				}

				// Job doesnt exist, create it.
				log.Info("forgetting the restic snapshot")
				err := r.DeleteResticSnapshot(backup)
				return reconcile.Result{RequeueAfter: 5 * time.Second}, err
			}

			removeFinalizer := false
			for _, condition := range newJob.Status.Conditions {
				if condition.Type == batchv1.JobComplete {
					log.Info("restic forget job complete")
					removeFinalizer = true
				} else if condition.Type == batchv1.JobFailed {
					log.Error("restic forget job failed")
					removeFinalizer = true
				}
			}
			if removeFinalizer {
				log.Infof("removing finalizer %s", Finalizer)

				// Remove this finalizer from the list and update it.
				backup.ObjectMeta.Finalizers = slice.Remove(backup.ObjectMeta.Finalizers, Finalizer)
				err := r.Update(context.Background(), backup)
				return reconcile.Result{}, err
			} else {
				log.Info("doing another loop")
				return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
			}
		}

		return reconcile.Result{}, nil
	}

	// Backup has completed or failed, return early.
	if backup.Status.Phase == v1.PhaseCompleted || backup.Status.Phase == v1.PhaseFailed {
		return reconcile.Result{}, nil
	}

	if _, found := backup.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the backup failed
		log.Info(fmt.Sprintf("Backup %s doesn't have a site label, skipping.", backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	return r.SyncJob(log, backup)
}

// SyncJob creates or updates the restic backup jobs.
func (r *ReconcileBackup) SyncJob(log log.Logger, backup *extensionv1.Backup) (reconcile.Result, error) {
	// Backup has completed or failed, return early.
	if backup.Status.Phase == v1.PhaseCompleted || backup.Status.Phase == v1.PhaseFailed {
		return reconcile.Result{}, nil
	}

	if _, found := backup.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the backup failed
		log.Info(fmt.Sprintf("Backup %s doesn't have a site label, skipping.", backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}
	var params = resticutils.PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "docker.io/restic/restic:0.9.5",
		MySQLImage:  "skpr/mtk-mysql",
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

// DeleteResticSnapshot creates the job to forget a restic snapshot.
func (r *ReconcileBackup) DeleteResticSnapshot(backup *extensionv1.Backup) error {
	var params = resticutils.PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "docker.io/restic/restic:0.9.5",
		MySQLImage:  "skpr/mtk-mysql",
		WorkingDir:  "/home/shepherd",
		Tags:        []string{},
	}

	spec, err := resticutils.PodSpecDelete(
		backup.Status.ResticID,
		backup.ObjectMeta.Namespace,
		backup.ObjectMeta.GetLabels()["site"],
		params,
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
			Name:      fmt.Sprintf("%s-delete-%s", resticutils.Prefix, backup.Name),
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

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.Client, job, func(obj runtime.Object) error {
		return nil
	})

	return err
}

// getResticIdFromJob parses output from a job's pods and returns a restic ID from the logs.
func getResticIdFromJob(c *rest.Config, job *batchv1.Job) (string, error) {
	kubeset, err := kubernetes.NewForConfig(c)
	var resticId string

	if err != nil {
		return resticId, err
	}

	pods, err := kubeset.CoreV1().Pods(job.ObjectMeta.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(job.Spec.Template.ObjectMeta.Labels)).String(),
	})
	if err != nil {
		return resticId, err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodSucceeded {
			continue
		}
		podLogs, err := getPodLogs(kubeset, pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
		if err != nil {
			return resticId, err
		}
		resticId = resticutils.ParseSnapshotID(podLogs)
		if resticId != "" {
			return resticId, nil
		}
	}

	return resticId, nil
}

// getPodLogs gets the logs from the restic container from a pod as a string.
func getPodLogs(kubeset *kubernetes.Clientset, namespace string, podName string) (string, error) {
	body, err := kubeset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: resticutils.ResticBackupContainerName,
	}).Stream()
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
