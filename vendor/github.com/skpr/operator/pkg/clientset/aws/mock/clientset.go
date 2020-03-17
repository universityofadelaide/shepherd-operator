package mock

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	clientset "github.com/skpr/operator/pkg/clientset/aws"
	"github.com/skpr/operator/pkg/clientset/aws/certificate"
	certificatemock "github.com/skpr/operator/pkg/clientset/aws/certificate/mock"
	"github.com/skpr/operator/pkg/clientset/aws/certificaterequest"
	certificaterequestmock "github.com/skpr/operator/pkg/clientset/aws/certificaterequest/mock"
	"github.com/skpr/operator/pkg/clientset/aws/cloudfront"
	cloudfrontmock "github.com/skpr/operator/pkg/clientset/aws/cloudfront/mock"
	"github.com/skpr/operator/pkg/clientset/aws/cloudfrontinvalidation"
	cloudfrontinvalidationmock "github.com/skpr/operator/pkg/clientset/aws/cloudfrontinvalidation/mock"
)

// Clientset used for mocking.
type Clientset struct {
	certificates            []*awsv1beta1.Certificate
	certificaterequests     []*awsv1beta1.CertificateRequest
	cloudfronts             []*awsv1beta1.CloudFront
	cloudfrontinvalidations []*awsv1beta1.CloudFrontInvalidation
}

// New clientset for interacting with Workflow objects.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	clientset := &Clientset{}

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Kind {
		case certificate.Kind:
			clientset.certificates = append(clientset.certificates, object.(*awsv1beta1.Certificate))
		case certificaterequest.Kind:
			clientset.certificaterequests = append(clientset.certificaterequests, object.(*awsv1beta1.CertificateRequest))
		case cloudfront.Kind:
			clientset.cloudfronts = append(clientset.cloudfronts, object.(*awsv1beta1.CloudFront))
		case cloudfrontinvalidation.Kind:
			clientset.cloudfrontinvalidations = append(clientset.cloudfrontinvalidations, object.(*awsv1beta1.CloudFrontInvalidation))
		default:
			return nil, fmt.Errorf("cannot find client for: %s", gvk.Kind)
		}
	}

	return clientset, nil
}

// Certificates within a namespace.
func (c *Clientset) Certificates(namespace string) certificate.Interface {
	filter := func(list []*awsv1beta1.Certificate, namespace string) []*awsv1beta1.Certificate {
		var filtered []*awsv1beta1.Certificate

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &certificatemock.Client{
		Namespace: namespace,
		Objects:   filter(c.certificates, namespace),
	}
}

// CertificateRequests within a namespace.
func (c *Clientset) CertificateRequests(namespace string) certificaterequest.Interface {
	filter := func(list []*awsv1beta1.CertificateRequest, namespace string) []*awsv1beta1.CertificateRequest {
		var filtered []*awsv1beta1.CertificateRequest

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &certificaterequestmock.Client{
		Namespace: namespace,
		Objects:   filter(c.certificaterequests, namespace),
	}
}

// CloudFronts within a namespace.
func (c *Clientset) CloudFronts(namespace string) cloudfront.Interface {
	filter := func(list []*awsv1beta1.CloudFront, namespace string) []*awsv1beta1.CloudFront {
		var filtered []*awsv1beta1.CloudFront

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &cloudfrontmock.Client{
		Namespace: namespace,
		Objects:   filter(c.cloudfronts, namespace),
	}
}

// CloudFrontInvalidations within a namespace.
func (c *Clientset) CloudFrontInvalidations(namespace string) cloudfrontinvalidation.Interface {
	filter := func(list []*awsv1beta1.CloudFrontInvalidation, namespace string) []*awsv1beta1.CloudFrontInvalidation {
		var filtered []*awsv1beta1.CloudFrontInvalidation

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &cloudfrontinvalidationmock.Client{
		Namespace: namespace,
		Objects:   filter(c.cloudfrontinvalidations, namespace),
	}
}
