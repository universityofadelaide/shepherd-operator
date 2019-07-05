package certificatevalidate

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
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
)

// ControllerName is used to identify this controller.
const ControllerName = "certificate-validate-controller"

// Add creates a new CertificateRequest Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, route53 route53iface.Route53API, params Params) error {
	return add(mgr, newReconciler(mgr, route53, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, route53 route53iface.Route53API, params Params) reconcile.Reconciler {
	return &ReconcileCertificateRequest{
		Client:   mgr.GetClient(),
		recorder: mgr.GetRecorder(ControllerName),
		scheme:   mgr.GetScheme(),
		route53:  route53,
		params:   params,
	}
}

// add a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Certificate
	return c.Watch(&source.Kind{Type: &awsv1beta1.CertificateRequest{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileCertificateRequest{}

// ReconcileCertificateRequest reconciles a Certificate object
type ReconcileCertificateRequest struct {
	client.Client
	route53  route53iface.Route53API
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	params   Params
}

// Params which are associated with this controller.
type Params struct {
	Zone   string
	Domain string
	TTL    int64
}

// Reconcile reads that state of the cluster for a Certificate object and makes changes based on the state read
// and what is in the Certificate.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=aws.skpr.io,resources=certificaterequests,verbs=get;list;watch;update;patch
func (r *ReconcileCertificateRequest) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconciler loop")

	// Fetch the CertificateRequest instance
	instance := &awsv1beta1.CertificateRequest{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	var changes []*route53.Change

	for _, validate := range instance.Status.Validate {
		if !strings.HasSuffix(validate.Name, r.params.Domain) {
			log.Infof("Skipping: %s", validate.Name)
			continue
		}

		log.Infof("Found: %s", validate.Name)

		changes = append(changes, &route53.Change{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: aws.String(validate.Name),
				Type: aws.String(validate.Type),
				ResourceRecords: []*route53.ResourceRecord{
					{
						Value: aws.String(validate.Value),
					},
				},
				TTL: aws.Int64(r.params.TTL),
			},
		})
	}

	if len(changes) > 0 {
		log.Info("Submitting change record request")

		_, err := r.route53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &route53.ChangeBatch{
				Changes: changes,
				Comment: aws.String("Auto created by github.com/skpr/operator"),
			},
			HostedZoneId: aws.String(r.params.Zone),
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	log.Info("Reconcile loop finished")

	return reconcile.Result{RequeueAfter: time.Duration(time.Minute * 5)}, nil
}
