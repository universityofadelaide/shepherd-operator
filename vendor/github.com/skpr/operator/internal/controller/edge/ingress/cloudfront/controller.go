package cloudfront

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	k8ssync "github.com/skpr/operator/pkg/utils/k8s/sync"
)

// ControllerName used for identifying which controller is performing an operation.
const ControllerName = "ingress-cloudfront-controller"

// Add creates a new Ingress Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params Params) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params Params) reconcile.Reconciler {
	return &ReconcileIngress{
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

	err = c.Watch(&source.Kind{Type: &extensionsv1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &edgev1beta1.Ingress{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &awsv1beta1.CloudFront{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &edgev1beta1.Ingress{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &awsv1beta1.Certificate{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &edgev1beta1.Ingress{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to Ingress
	return c.Watch(&source.Kind{Type: &edgev1beta1.Ingress{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileIngress{}

// ReconcileIngress reconciles a Ingress object
type ReconcileIngress struct {
	client.Client
	scheme *runtime.Scheme
	params Params
}

// Params which inform the Reconciler.
type Params struct {
	OriginEndpoint string
	OriginPolicy   string
	OriginTimeout  int64
}

// Reconcile reads that state of the cluster for a Ingress object and makes changes based on the state read
// and what is in the Ingress.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=edge.skpr.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=edge.skpr.io,resources=ingresses/status,verbs=get;update;patch
func (r *ReconcileIngress) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting Reconcile Loop")

	ingress := &edgev1beta1.Ingress{}

	err := r.Get(context.TODO(), request.NamespacedName, ingress)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	log.Info("Syncing Ingress")

	status, err := r.Sync(log, ingress)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "sync status failed")
	}

	status.ObservedGeneration = ingress.ObjectMeta.Generation

	err = r.SyncStatus(ingress, status)
	if err != nil {
		log.Error(err, "Status status failed")
		return reconcile.Result{}, errors.Wrap(err, "sync status failed")
	}

	log.Info("Reconcile Loop Finished")

	return reconcile.Result{}, nil
}

// Sync Ingress / Certificate / CloudFront resources.
func (r *ReconcileIngress) Sync(log log.Logger, instance *edgev1beta1.Ingress) (edgev1beta1.IngressStatus, error) {
	var status edgev1beta1.IngressStatus

	ingress := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.ObjectMeta.Name,
			Namespace: instance.ObjectMeta.Namespace,
		},
	}

	for _, route := range append(instance.Spec.Routes.Secondary, instance.Spec.Routes.Primary) {
		rule := extensionsv1beta1.IngressRule{
			Host: route.Domain,
			IngressRuleValue: extensionsv1beta1.IngressRuleValue{
				HTTP: &extensionsv1beta1.HTTPIngressRuleValue{},
			},
		}

		for _, path := range route.Subpaths {
			rule.IngressRuleValue.HTTP.Paths = append(rule.IngressRuleValue.HTTP.Paths, extensionsv1beta1.HTTPIngressPath{
				Path: path,
				Backend: extensionsv1beta1.IngressBackend{
					ServiceName: instance.Spec.Service.Name,
					ServicePort: intstr.FromInt(instance.Spec.Service.Port),
				},
			})
		}

		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, ingress, k8ssync.Ingress(instance, *ingress, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Ingress")
	}
	log.Infof("Synced Ingress %s with status: %s", ingress.ObjectMeta.Name, result)

	certificate := &awsv1beta1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.ObjectMeta.Name,
			Namespace: instance.ObjectMeta.Namespace,
		},
		Spec: awsv1beta1.CertificateSpec{
			Request: awsv1beta1.CertificateRequestSpec{
				CommonName:     instance.Spec.Routes.Primary.Domain,
				AlternateNames: getDomainsFromRoutes(instance.Spec.Routes.Secondary),
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, certificate, k8ssync.Certificate(instance, certificate.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Certificate")
	}

	status.Certificate = edgev1beta1.IngressStatusCertificateRef{
		Name:    certificate.ObjectMeta.Name,
		Details: certificate.Status,
	}

	log.Infof("Synced Certificate %s with status: %s", certificate.ObjectMeta.Name, result)

	cloudfront := &awsv1beta1.CloudFront{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.ObjectMeta.Name,
			Namespace: instance.ObjectMeta.Namespace,
		},
		Spec: awsv1beta1.CloudFrontSpec{
			Aliases: certificate.Status.Active.Details.Domains,
			Certificate: awsv1beta1.CloudFrontSpecCertificate{
				ARN: certificate.Status.Active.Details.ARN,
			},
			Behavior: awsv1beta1.CloudFrontSpecBehavior{
				Whitelist: instance.Spec.Whitelist,
			},
			Origin: awsv1beta1.CloudFrontSpecOrigin{
				Endpoint: r.params.OriginEndpoint,
				Policy:   r.params.OriginPolicy,
				Timeout:  r.params.OriginTimeout,
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, cloudfront, k8ssync.CloudFront(instance, cloudfront.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CloudFront")
	}

	status.CloudFront = edgev1beta1.IngressStatusCloudFrontRef{
		Name:    cloudfront.ObjectMeta.Name,
		Details: cloudfront.Status,
	}

	log.Infof("Synced CloudFront %s with status: %s", cloudfront.ObjectMeta.Name, result)

	return status, nil
}

// SyncStatus with the Ingress object.
func (r *ReconcileIngress) SyncStatus(ingress *edgev1beta1.Ingress, status edgev1beta1.IngressStatus) error {
	if diff := deep.Equal(ingress.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		ingress.Status = status

		return r.Status().Update(context.TODO(), ingress)
	}

	return nil
}
