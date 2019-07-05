package certificate

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/k8s/events"
	k8ssync "github.com/skpr/operator/pkg/utils/k8s/sync"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

const (
	// ControllerName is used to identify this controller.
	ControllerName = "certificate-controller"
	// LabelCertificate for creating and querying certificates.
	LabelCertificate = "certificate"
)

// Add creates a new Certificate Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, params Params) error {
	return add(mgr, newReconciler(mgr, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, params Params) reconcile.Reconciler {
	return &ReconcileCertificate{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetRecorder(ControllerName),
		params:   params,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &awsv1beta1.CertificateRequest{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &awsv1beta1.Certificate{},
	})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &awsv1beta1.Certificate{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileCertificate{}

// ReconcileCertificate reconciles a Certificate object
type ReconcileCertificate struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
	params   Params
}

// Params for this controller.
type Params struct {
	Retention int
}

// Reconcile reads that state of the cluster for a Certificate object and makes changes based on the state read
// and what is in the Certificate.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=aws.skpr.io,resources=certificate,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aws.skpr.io,resources=certificaterequests,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileCertificate) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconciler loop")

	certificate := &awsv1beta1.Certificate{}

	err := r.Get(context.TODO(), request.NamespacedName, certificate)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, errors.Wrap(err, "failed to get Certificate")
	}

	log.Info("Syncing CertificateRequest")

	status, err := r.Sync(certificate)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to sync CertificateRequest")
	}

	log.Info("Syncing Status")

	err = r.SyncStatus(certificate, status)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to sync CertificateRequest")
	}

	log.Info("Cleaning up old CertificateRequests")

	err = r.Cleanup(certificate, r.params.Retention)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to cleanup old CertificateRequests")
	}

	log.Info("Finished reconcile loop")

	return reconcile.Result{}, nil
}

// Sync will create and/or return an existing CertificateRequest.
func (r *ReconcileCertificate) Sync(certificate *awsv1beta1.Certificate) (awsv1beta1.CertificateStatus, error) {
	status := awsv1beta1.CertificateStatus{
		ObservedGeneration: certificate.ObjectMeta.Generation,
	}

	selector := map[string]string{
		LabelCertificate: certificate.ObjectMeta.Name,
	}

	// This is the CertificateRequest we want based on the CommonName and AlternativeNames provided.
	desired := &awsv1beta1.CertificateRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", certificate.ObjectMeta.Name, certificate.Spec.Request.Hash()),
			Namespace: certificate.ObjectMeta.Namespace,
			Labels:    selector,
		},
		Spec: certificate.Spec.Request,
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, desired, k8ssync.CertificateRequest(certificate, desired.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CertificateRequest")
	}

	list := &awsv1beta1.CertificateRequestList{}

	query := &client.ListOptions{
		Namespace:     certificate.ObjectMeta.Namespace,
		LabelSelector: labels.SelectorFromSet(labels.Set(selector)),
	}

	err = r.List(context.TODO(), query, list)
	if err != nil {
		return status, err
	}

	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].ObjectMeta.CreationTimestamp.Time.Unix() < list.Items[j].ObjectMeta.CreationTimestamp.Time.Unix()
	})

	status.Desired = requestToReference(*desired)
	status.Active = getActiveRequestStatus(*desired, list.Items)

	for _, item := range list.Items {
		status.Requests = append(status.Requests, requestToReference(item))
	}

	return status, nil
}

// SyncStatus back to Certificate.
func (r *ReconcileCertificate) SyncStatus(certificate *awsv1beta1.Certificate, status awsv1beta1.CertificateStatus) error {
	if diff := deep.Equal(certificate.Status, status); diff != nil {
		if certificate.Status.Desired.Name != status.Desired.Name {
			r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventCreate, "Creating CertificateRequest: %s", status.Desired.Name)
		}

		if certificate.Status.Active.Name != status.Active.Name {
			r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventCreate, "Promoting CertificateRequest to active: %s", status.Active.Name)
		}

		certificate.Status = status

		return r.Status().Update(context.TODO(), certificate)
	}

	return nil
}

// Cleanup old Requests which are not issued.
func (r *ReconcileCertificate) Cleanup(certificate *awsv1beta1.Certificate, retention int) error {
	for key, item := range certificate.Status.Requests {
		if key < retention {
			continue
		}

		if item.Name == certificate.Status.Active.Name {
			continue
		}

		r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventCreate, "Cleaning up old CertificateRequest: %s", item.Name)

		delete := &awsv1beta1.CertificateRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      item.Name,
				Namespace: certificate.ObjectMeta.Namespace,
			},
		}

		err := r.Delete(context.TODO(), delete)
		if err != nil {
			return err
		}
	}

	return nil
}
