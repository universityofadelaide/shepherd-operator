package backup

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	v1 "gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/apis/meta/v1"
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

	extensionv1 "gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/apis/extension/v1"
	"gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/utils/k8s/sync"
	resticutils "gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/utils/restic"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "backup-restic-controller"

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
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
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
