package image

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	jobutils "github.com/skpr/operator/pkg/utils/k8s/job"
	"github.com/skpr/operator/pkg/utils/k8s/sync"
	mysqlimage "github.com/skpr/operator/pkg/utils/mysql/image"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "mysql-image"

// Add creates a new Image Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params mysqlimage.GenerateParams) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params mysqlimage.GenerateParams) reconcile.Reconciler {
	return &ReconcileImage{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		params: params,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mysqlv1beta1.Image{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to Image
	return c.Watch(&source.Kind{Type: &mysqlv1beta1.Image{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileImage{}

// ReconcileImage reconciles a Image object
type ReconcileImage struct {
	client.Client
	scheme *runtime.Scheme
	params mysqlimage.GenerateParams
}

// Reconcile reads that state of the cluster for a Image object and makes changes based on the state read
// and what is in the Image.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.skpr.io,resources=images,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileImage) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	image := &mysqlv1beta1.Image{}

	err := r.Get(context.TODO(), request.NamespacedName, image)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	pod, configmap, err := mysqlimage.Generate(image, r.params)
	if err != nil {
		return reconcile.Result{}, err
	}

	job, err := jobutils.NewFromPod(pod.ObjectMeta, pod.Spec)
	if err != nil {
		return reconcile.Result{}, err
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, job, sync.Job(image, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced Job %s with status: %s", job.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, configmap, sync.ConfigMap(image, configmap.Data, configmap.BinaryData, true, r.scheme))
	if err != nil {
		return reconcile.Result{}, err
	}
	log.Infof("Synced ConfigMap %s with status: %s", configmap.ObjectMeta.Name, result)

	status := mysqlv1beta1.ImageStatus{
		Phase:          jobutils.GetPhase(job.Status),
		StartTime:      job.Status.StartTime,
		CompletionTime: job.Status.CompletionTime,
	}

	if diff := deep.Equal(image.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		image.Status = status

		err := r.Status().Update(context.TODO(), image)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	return reconcile.Result{}, nil
}
