package cloudfrontinvalidation

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "cloudfront-invalidation-controller"
)

// Add creates a new CloudFrontInvalidation Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, cloudfront cloudfrontiface.CloudFrontAPI) error {
	return add(mgr, newReconciler(mgr, cloudfront))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, cloudfront cloudfrontiface.CloudFrontAPI) reconcile.Reconciler {
	return &ReconcileCloudFrontInvalidation{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		recorder:   mgr.GetRecorder(ControllerName),
		cloudfront: cloudfront,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to CloudFrontInvalidation
	return c.Watch(&source.Kind{Type: &awsv1beta1.CloudFrontInvalidation{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileCloudFrontInvalidation{}

// ReconcileCloudFrontInvalidation reconciles a CloudFrontInvalidation object
type ReconcileCloudFrontInvalidation struct {
	client.Client
	scheme     *runtime.Scheme
	recorder   record.EventRecorder
	cloudfront cloudfrontiface.CloudFrontAPI
}

// Reconcile reads that state of the cluster for a CloudFrontInvalidation object and makes changes based on the state read
// and what is in the CloudFrontInvalidation.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=aws.skpr.io,resources=cloudfrontinvalidations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws.skpr.io,resources=cloudfrontinvalidations/status,verbs=get;update;patch
func (r *ReconcileCloudFrontInvalidation) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	invalidation := &awsv1beta1.CloudFrontInvalidation{}

	err := r.Get(context.TODO(), request.NamespacedName, invalidation)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	distribution := &awsv1beta1.CloudFront{}

	query := types.NamespacedName{
		Name:      invalidation.Spec.Distribution,
		Namespace: invalidation.ObjectMeta.Namespace,
	}

	err = r.Get(context.TODO(), query, distribution)
	if err != nil {
		return reconcile.Result{}, err
	}

	var status awsv1beta1.CloudFrontInvalidationStatus

	if invalidation.Status.ID != "" {
		status, err = r.DescribeInvalidation(distribution, invalidation)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		status, err = r.CreateInvalidation(distribution, invalidation)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	status.ObservedGeneration = invalidation.ObjectMeta.Generation

	if diff := deep.Equal(invalidation.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		invalidation.Status = status

		err := r.Status().Update(context.TODO(), invalidation)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	// Requeue if we are still waiting for the invalidation to finish.
	if invalidation.Status.State != awsv1beta1.CloudFrontInvalidationCompleted {
		log.Info("Reconcile loop finished, requeuing at a frequent interval while waiting for invalidation to finish")

		return reconcile.Result{RequeueAfter: time.Duration(time.Second * 15)}, nil
	}

	return reconcile.Result{}, nil
}

// CreateInvalidation will process an invalidation request assigned to a distribution.
func (r *ReconcileCloudFrontInvalidation) CreateInvalidation(distribution *awsv1beta1.CloudFront, invalidation *awsv1beta1.CloudFrontInvalidation) (awsv1beta1.CloudFrontInvalidationStatus, error) {
	var status awsv1beta1.CloudFrontInvalidationStatus

	resp, err := r.cloudfront.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distribution.Status.ID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().String()),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(int64(len(invalidation.Spec.Paths))),
				Items:    aws.StringSlice(invalidation.Spec.Paths),
			},
		},
	})
	if err != nil {
		return status, err
	}

	status.ID = *resp.Invalidation.Id
	status.Created = resp.Invalidation.CreateTime.String()
	status.State = *resp.Invalidation.Status

	return status, nil
}

// DescribeInvalidation will lookup an invalidation request assigned to a distribution.
func (r *ReconcileCloudFrontInvalidation) DescribeInvalidation(distribution *awsv1beta1.CloudFront, invalidation *awsv1beta1.CloudFrontInvalidation) (awsv1beta1.CloudFrontInvalidationStatus, error) {
	var status awsv1beta1.CloudFrontInvalidationStatus

	resp, err := r.cloudfront.GetInvalidation(&cloudfront.GetInvalidationInput{
		DistributionId: aws.String(distribution.Status.ID),
		Id:             aws.String(invalidation.Status.ID),
	})
	if err != nil {
		return status, err
	}

	status.ID = *resp.Invalidation.Id
	status.Created = resp.Invalidation.CreateTime.String()
	status.State = *resp.Invalidation.Status

	return status, nil
}
