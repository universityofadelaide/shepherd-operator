package certificate

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Kind of object this client handles.
const Kind = "Certificate"

// Interface declares interactions with the Certificate objects.
type Interface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CertificateList, error)
	Get(name string, opts metav1.GetOptions) (*awsv1beta1.Certificate, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*awsv1beta1.Certificate) (*awsv1beta1.Certificate, error)
	Update(*awsv1beta1.Certificate) (*awsv1beta1.Certificate, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Certificate objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Certificates.
func (c *Client) List(opts metav1.ListOptions) (*awsv1beta1.CertificateList, error) {
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

// Get a Certificate.
func (c *Client) Get(name string, opts metav1.GetOptions) (*awsv1beta1.Certificate, error) {
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

// Exists checks if a Certificate is present.
func (c *Client) Exists(name string, opts metav1.GetOptions) (bool, error) {
	_, err := c.Get(name, opts)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Create a Certificate.
func (c *Client) Create(certificate *awsv1beta1.Certificate) (*awsv1beta1.Certificate, error) {
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

// Update a Certificate.
func (c *Client) Update(certificate *awsv1beta1.Certificate) (*awsv1beta1.Certificate, error) {
	result := awsv1beta1.Certificate{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("certificates").
		Name(certificate.Name).
		Body(certificate).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Certificate.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("certificates").
		Name(name).
		Do().
		Error()
}

// Watch for Certificates.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
