package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// CertificateInterface declares interactions with the Certificate objects.
type CertificateInterface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CertificateList, error)
	Get(name string, options metav1.GetOptions) (*awsv1beta1.Certificate, error)
	Create(*awsv1beta1.Certificate) (*awsv1beta1.Certificate, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type certificateClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Certificates.
func (c *certificateClient) List(opts metav1.ListOptions) (*awsv1beta1.CertificateList, error) {
	result := awsv1beta1.CertificateList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Certificate.
func (c *certificateClient) Get(name string, opts metav1.GetOptions) (*awsv1beta1.Certificate, error) {
	result := awsv1beta1.Certificate{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificates").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Certificate.
func (c *certificateClient) Create(certificate *awsv1beta1.Certificate) (*awsv1beta1.Certificate, error) {
	result := awsv1beta1.Certificate{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("certificates").
		Body(certificate).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Certificates.
func (c *certificateClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
