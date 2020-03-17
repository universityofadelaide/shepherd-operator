package certificaterequest

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/skpr/operator/pkg/utils/uid"
	corev1 "k8s.io/api/core/v1"
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
	"github.com/skpr/operator/pkg/utils/k8s/events"
	"github.com/skpr/operator/pkg/utils/slice"
)

const (
	// FinalizerName for tagging objects for action on deletion (delete remote source).
	FinalizerName = "certificaterequest.aws.skpr.io"
	// ControllerName is used to identify this controller.
	ControllerName = "certificate-request-controller"
)

// Add creates a new CertificateRequest Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, acm acmiface.ACMAPI, params Params) error {
	return add(mgr, newReconciler(mgr, acm, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, acm acmiface.ACMAPI, params Params) reconcile.Reconciler {
	return &ReconcileCertificateRequest{
		Client:   mgr.GetClient(),
		recorder: mgr.GetRecorder(ControllerName),
		scheme:   mgr.GetScheme(),
		acm:      acm,
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
	acm      acmiface.ACMAPI
	recorder record.EventRecorder
	scheme   *runtime.Scheme
	params   Params
}

// Params for this controller.
type Params struct {
	Prefix string
}

// Reconcile reads that state of the cluster for a Certificate object and makes changes based on the state read
// and what is in the Certificate.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=aws.skpr.io,resources=certificates,verbs=get;list;watch;update;patch
func (r *ReconcileCertificateRequest) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	certificate := &awsv1beta1.CertificateRequest{}

	err := r.Get(context.TODO(), request.NamespacedName, certificate)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if certificate.ObjectMeta.DeletionTimestamp.IsZero() {
		// We need to attach a finalizer to this object so we can delete the certificate while
		// we are also removing this object.
		if !slice.Contains(certificate.ObjectMeta.Finalizers, FinalizerName) {
			log.Info("Adding finalizer:", FinalizerName)

			certificate.ObjectMeta.Finalizers = append(certificate.ObjectMeta.Finalizers, FinalizerName)
			if err := r.Update(context.Background(), certificate); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		if slice.Contains(certificate.ObjectMeta.Finalizers, FinalizerName) {
			log.Info("Deleting certificate")

			err := r.DeleteCertificate(certificate)
			if err != nil {
				return reconcile.Result{}, err
			}

			// We can now remove the finalizer to let this object be removed.
			certificate.ObjectMeta.Finalizers = slice.Remove(certificate.ObjectMeta.Finalizers, FinalizerName)
			if err := r.Update(context.Background(), certificate); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	}

	arn, err := r.CreateCertificate(certificate)
	if err != nil {
		return reconcile.Result{}, err
	}

	status, err := r.DescribeCertificate(arn)
	if err != nil {
		return reconcile.Result{}, err
	}

	status.ObservedGeneration = certificate.ObjectMeta.Generation

	if diff := deep.Equal(certificate.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		certificate.Status = status

		err := r.Status().Update(context.TODO(), certificate)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	// Requeue if we are still waiting for the certificate to become available.
	if certificate.Status.State != acm.CertificateStatusIssued {
		log.Info("Reconcile loop finished, requeuing at a frequent interval while waiting for provisioning to finish")

		return reconcile.Result{RequeueAfter: time.Duration(time.Second * 30)}, nil
	}

	log.Info("Reconcile loop finished")

	return reconcile.Result{RequeueAfter: time.Duration(time.Minute * 5)}, nil
}

// CreateCertificate as part of a request.
func (r *ReconcileCertificateRequest) CreateCertificate(certificate *awsv1beta1.CertificateRequest) (string, error) {
	if certificate.Status.ARN != "" {
		return certificate.Status.ARN, nil
	}

	r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventCreate, "Requesting certificate")

	token, err := uid.GetToken(certificate.ObjectMeta.UID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}

	input := &acm.RequestCertificateInput{
		IdempotencyToken: aws.String(token),
		DomainName:       aws.String(certificate.Spec.CommonName),
		ValidationMethod: aws.String(acm.ValidationMethodDns),
	}

	if len(certificate.Spec.AlternateNames) > 0 {
		input.SubjectAlternativeNames = aws.StringSlice(certificate.Spec.AlternateNames)
	}

	resp, err := r.acm.RequestCertificate(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to request certificate")
	}

	r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventCreate, "Request created: %s", *resp.CertificateArn)

	// We want to requeue this object so we can let the reconcile loop describe the object now it's created.
	return *resp.CertificateArn, nil
}

// DescribeCertificate to get the latest status.
func (r *ReconcileCertificateRequest) DescribeCertificate(arn string) (awsv1beta1.CertificateRequestStatus, error) {
	var status awsv1beta1.CertificateRequestStatus

	resp, err := r.acm.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	})
	if err != nil {
		return status, errors.Wrap(err, "failed to describe certificate")
	}

	status.ARN = *resp.Certificate.CertificateArn
	status.State = *resp.Certificate.Status

	// The SubjectAlternativeNames include the primary domain.
	// https://docs.aws.amazon.com/acm/latest/APIReference/API_CertificateDetail.html
	status.Domains = aws.StringValueSlice(resp.Certificate.SubjectAlternativeNames)

	for _, val := range resp.Certificate.DomainValidationOptions {
		status.Validate = append(status.Validate, awsv1beta1.ValidateRecord{
			Name:   *val.ResourceRecord.Name,
			Type:   *val.ResourceRecord.Type,
			Status: *val.ValidationStatus,
			Value:  *val.ResourceRecord.Value,
		})
	}

	return status, nil
}

// DeleteCertificate as part of cleanup.
func (r *ReconcileCertificateRequest) DeleteCertificate(certificate *awsv1beta1.CertificateRequest) error {
	r.recorder.Eventf(certificate, corev1.EventTypeNormal, events.EventDelete, "Deleting certificate")

	if certificate.Status.ARN == "" {
		return nil
	}

	_, err := r.acm.DeleteCertificate(&acm.DeleteCertificateInput{
		CertificateArn: aws.String(certificate.Status.ARN),
	})
	if err != nil && !kerrors.IsNotFound(err) {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == acm.ErrCodeResourceNotFoundException {
				return nil
			}
		}

		return errors.Wrap(err, "failed to delete certificate")
	}

	return nil
}
