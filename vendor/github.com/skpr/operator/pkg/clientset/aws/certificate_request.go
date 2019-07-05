package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// CertificateRequestInterface declares interactions with the CertificateRequest objects.
type CertificateRequestInterface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CertificateRequestList, error)
	Get(name string, options metav1.GetOptions) (*awsv1beta1.CertificateRequest, error)
	Create(*awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type certificateRequestClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all CertificateRequests.
func (c *certificateRequestClient) List(opts metav1.ListOptions) (*awsv1beta1.CertificateRequestList, error) {
	result := awsv1beta1.CertificateRequestList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an CertificateRequest.
func (c *certificateRequestClient) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CertificateRequest, error) {
	result := awsv1beta1.CertificateRequest{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an CertificateRequest.
func (c *certificateRequestClient) Create(certificaterequest *awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error) {
	result := awsv1beta1.CertificateRequest{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		Body(certificaterequest).
		Do().
		Into(&result)

	return &result, err
}

// Watch for CertificateRequests.
func (c *certificateRequestClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
