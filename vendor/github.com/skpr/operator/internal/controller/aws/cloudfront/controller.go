package cloudfront

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/slice"
)

const (
	// Finalizer used to trigger a deletion of the user prior to the object being deleted.
	Finalizer = "cloudfronts.aws.skpr.io"
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "cloudfront-controller"
)

// Add creates a new CloudFront Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, cloudfront cloudfrontiface.CloudFrontAPI, params Params) error {
	return add(mgr, newReconciler(mgr, cloudfront, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, cloudfront cloudfrontiface.CloudFrontAPI, params Params) reconcile.Reconciler {
	return &ReconcileCloudFront{
		Client:     mgr.GetClient(),
		cloudfront: cloudfront,
		recorder:   mgr.GetRecorder(ControllerName),
		scheme:     mgr.GetScheme(),
		params:     params,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to CloudFront
	return c.Watch(&source.Kind{Type: &awsv1beta1.CloudFront{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileCloudFront{}

// ReconcileCloudFront reconciles a CloudFront object
type ReconcileCloudFront struct {
	client.Client
	cloudfront cloudfrontiface.CloudFrontAPI
	recorder   record.EventRecorder
	scheme     *runtime.Scheme
	params     Params
}

// Params which are passed into this reconciler.
type Params struct {
	Prefix        string
	LoggingBucket string
}

// Reconcile reads that state of the cluster for a CloudFront object and makes changes based on the state read
// and what is in the CloudFront.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=aws.skpr.io,resources=cloudfronts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws.skpr.io,resources=cloudfronts/status,verbs=get;update;patch
func (r *ReconcileCloudFront) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	instance := &awsv1beta1.CloudFront{}

	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// https://book.kubebuilder.io/beyond_basics/using_finalizers.html
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !slice.Contains(instance.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Adding finalizer")

			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		// The object is being deleted, ensure that we have the finalizer and delete the CloudFront distribution.
		if slice.Contains(instance.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Deleting CloudFront distribution")

			// our finalizer is present, so lets handle our external dependency
			err := r.DeleteExternal(instance)
			if err != nil {
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = slice.Remove(instance.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), instance); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}

		return reconcile.Result{}, nil
	}

	log.Info("Syncing with CloudFront distribution")

	// Unique identifier to ensure we don't create more than one CloudFront per instance.
	generatedReference := fmt.Sprintf("%s-%s-%s", r.params.Prefix, instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)

	// When syncing we look to see if there is an existing reference which has been used.
	// This also allows for custom references to be used when importing an existing CloudFront distribution.
	if instance.Status.CallerReference == "" {
		instance.Status.CallerReference = generatedReference
	}

	distribution, err := r.SyncExternal(log, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	status := awsv1beta1.CloudFrontStatus{
		ObservedGeneration: instance.ObjectMeta.Generation,
		CallerReference:    instance.Status.CallerReference,
	}

	if distribution != nil {
		if distribution.Id != nil {
			status.ID = *distribution.Id
		}

		if distribution.Status != nil {
			status.State = *distribution.Status
		}

		if distribution.DomainName != nil {
			status.DomainName = *distribution.DomainName
		}

		if distribution.InProgressInvalidationBatches != nil {
			status.RunningInvalidations = *distribution.InProgressInvalidationBatches
		}
	}

	if diff := deep.Equal(instance.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		instance.Status = status

		err := r.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	if instance.Status.State == awsv1beta1.CloudFrontStateInProgress {
		log.Info("Reconcile loop finished, requeuing at a frequent interval while waiting for provisioning to finish")

		return reconcile.Result{RequeueAfter: time.Duration(time.Second * 15)}, nil
	}

	log.Info("Reconcile loop finished")

	return reconcile.Result{RequeueAfter: time.Duration(time.Minute)}, nil
}
