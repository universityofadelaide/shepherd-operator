package ses

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	"github.com/skpr/operator/pkg/utils/aws/policy"
	sesutils "github.com/skpr/operator/pkg/utils/aws/ses"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	"github.com/skpr/operator/pkg/utils/slice"
)

const (
	// Finalizer used to trigger a deletion of the user prior to the object being deleted.
	Finalizer = "ses.aws.skpr.io"
	// ControllerName is used to identify this controller.
	ControllerName = "ses-controller"
)

// Add creates a new SMTP Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, iam iamiface.IAMAPI, ses sesiface.SESAPI, params Params) error {
	return add(mgr, newReconciler(mgr, iam, ses, params))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, iam iamiface.IAMAPI, ses sesiface.SESAPI, params Params) reconcile.Reconciler {
	return &ReconcileSMTP{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		iam:    iam,
		ses:    ses,
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

	// Watch for changes to SMTP
	return c.Watch(&source.Kind{Type: &extensionsv1beta1.SMTP{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileSMTP{}

// ReconcileSMTP reconciles a SMTP object
type ReconcileSMTP struct {
	client.Client
	scheme *runtime.Scheme
	ses    sesiface.SESAPI
	iam    iamiface.IAMAPI
	params Params
}

// Params which are provided from external sources.
type Params struct {
	Prefix   string
	Hostname string
	Port     int
}

// Reconcile reads that state of the cluster for a SMTP object and makes changes based on the state read
// and what is in the SMTP.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=smtps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.skpr.io,resources=smtps/status,verbs=get;update;patch
func (r *ReconcileSMTP) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting reconcile loop")

	smtp := &extensionsv1beta1.SMTP{}

	err := r.Get(context.TODO(), request.NamespacedName, smtp)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	name := fmt.Sprintf("%s-%s-%s", r.params.Prefix, smtp.ObjectMeta.Namespace, smtp.ObjectMeta.Name)

	// https://book.kubebuilder.io/beyond_basics/using_finalizers.html
	if smtp.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object.
		if !slice.Contains(smtp.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Adding finalizer")

			smtp.ObjectMeta.Finalizers = append(smtp.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), smtp); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}
	} else {
		// The object is being deleted, ensure that we have the finalizer and delete the IAM user.
		if slice.Contains(smtp.ObjectMeta.Finalizers, Finalizer) {
			log.Info("Deleting IAM user")

			// our finalizer is present, so lets handle our external dependency
			err := r.DeleteUser(name)
			if err != nil {
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			smtp.ObjectMeta.Finalizers = slice.Remove(smtp.ObjectMeta.Finalizers, Finalizer)
			if err := r.Update(context.Background(), smtp); err != nil {
				return reconcile.Result{Requeue: true}, nil
			}
		}

		return reconcile.Result{}, nil
	}

	if smtp.Spec.From.Address == "" {
		log.Info("Skipping because .Spec.From.Address is not set")
		return reconcile.Result{}, nil
	}

	log.Info("Syncing email verification")

	address, err := r.SyncEmailVerification(smtp.Spec.From.Address)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, errors.Wrap(err, "failed to sync FROM address")
	}

	status := extensionsv1beta1.SMTPStatus{
		Connection: extensionsv1beta1.SMTPStatusConnection{
			Hostname: r.params.Hostname,
			Port:     r.params.Port,
		},
		Verification: extensionsv1beta1.SMTPStatusVerification{
			Address: address,
		},
	}

	err = r.SyncUser(name)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, errors.Wrap(err, "failed to sync user")
	}

	if smtp.Status.Connection.Username == "" || smtp.Status.Connection.Password == "" {
		log.Info("Syncing access keys")

		username, password, err := r.CreateAccessKeys(name)
		if err != nil {
			log.Error(err)
			return reconcile.Result{}, errors.Wrap(err, "failed to create credentials")
		}

		smtpPassword, err := sesutils.PasswordFromSecretKey(password)
		if err != nil {
			log.Error(err)
			return reconcile.Result{}, errors.Wrap(err, "failed to get password")
		}

		status.Connection.Username = username
		status.Connection.Password = smtpPassword
	} else {
		status.Connection.Username = smtp.Status.Connection.Username
		status.Connection.Password = smtp.Status.Connection.Password
	}

	log.Info("Syncing IAM policy")

	err = r.SyncPolicy(name, smtp.Spec.From.Address)
	if err != nil {
		log.Error(err)
		return reconcile.Result{}, errors.Wrap(err, "failed to sync policy")
	}

	if diff := deep.Equal(smtp.Status, status); diff != nil {
		log.Info(fmt.Sprintf("Status change dectected: %s", diff))

		smtp.Status = status

		err := r.Status().Update(context.TODO(), smtp)
		if err != nil {
			log.Error(err)
			return reconcile.Result{}, errors.Wrap(err, "failed to update status")
		}
	}

	if smtp.Status.Verification.Address == ses.CustomMailFromStatusPending {
		log.Info("Reconcile loop finished, requeuing at a frequent interval while waiting for provisioning to finish")

		return reconcile.Result{RequeueAfter: time.Duration(time.Second * 15)}, nil
	}

	log.Info("Reconcile loop finished")

	return reconcile.Result{RequeueAfter: time.Duration(time.Minute * 5)}, nil
}

// SyncUser with AWS IAM.
func (r *ReconcileSMTP) SyncUser(name string) error {
	_, err := r.iam.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == iam.ErrCodeEntityAlreadyExistsException {
				// The user has been created and we don't need to load it.
				return nil
			}
		}

		return err
	}

	return nil
}

// DeleteUser with AWS IAM.
func (r *ReconcileSMTP) DeleteUser(name string) error {
	_, err := r.iam.DeleteUser(&iam.DeleteUserInput{
		UserName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil
			}
		}

		return err
	}

	return nil
}

// CreateAccessKeys for accessing the AWS SES service.
func (r *ReconcileSMTP) CreateAccessKeys(name string) (string, string, error) {
	keys, err := r.iam.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(name),
	})
	if err != nil {
		return "", "", err
	}

	// Clear out old access keys. There can only be one.
	for _, key := range keys.AccessKeyMetadata {
		_, err := r.iam.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			UserName:    key.UserName,
			AccessKeyId: key.AccessKeyId,
		})
		if err != nil {
			return "", "", err
		}
	}

	key, err := r.iam.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(name),
	})
	if err != nil {
		return "", "", err
	}

	return *key.AccessKey.AccessKeyId, *key.AccessKey.SecretAccessKey, nil
}

// SyncPolicy to allow IAM user to send email.
func (r *ReconcileSMTP) SyncPolicy(name, address string) error {
	document, err := policy.Print(policy.Document{
		Version: "2012-10-17",
		Statement: []policy.Statement{
			{
				Effect: "Allow",
				Action: []string{
					"ses:SendEmail",
					"ses:SendRawEmail",
				},
				Resource: "*",
				Condition: policy.Conditions{
					StringEquals: map[string]string{
						"ses:FromAddress": address,
					},
				},
			},
		},
	})

	_, err = r.iam.PutUserPolicy(&iam.PutUserPolicyInput{
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(document),
		UserName:       aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// SyncEmailVerification to allow users to send email FROM an address.
func (r *ReconcileSMTP) SyncEmailVerification(address string) (string, error) {
	result, err := r.ses.GetIdentityVerificationAttributes(&ses.GetIdentityVerificationAttributesInput{
		Identities: []*string{
			aws.String(address),
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to list verification identities")
	}

	if val, ok := result.VerificationAttributes[address]; ok {
		return *val.VerificationStatus, nil
	}

	_, err = r.ses.VerifyEmailAddress(&ses.VerifyEmailAddressInput{
		EmailAddress: aws.String(address),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to submit verification request")
	}

	return r.SyncEmailVerification(address)
}
