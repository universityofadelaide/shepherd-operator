package certificaterequest

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Kind of object this client handles.
const Kind = "CertificateRequest"

// Interface declares interactions with the CertificateRequest objects.
type Interface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CertificateRequestList, error)
	Get(name string, opts metav1.GetOptions) (*awsv1beta1.CertificateRequest, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error)
	Update(*awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with CertificateRequest objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all CertificateRequests.
func (c *Client) List(opts metav1.ListOptions) (*awsv1beta1.CertificateRequestList, error) {
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

// Get a CertificateRequest.
func (c *Client) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CertificateRequest, error) {
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

// Exists checks if a CertificateRequest is present.
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

// Create a CertificateRequest.
func (c *Client) Create(certificaterequest *awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error) {
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

// Update a CertificateRequest.
func (c *Client) Update(certificaterequest *awsv1beta1.CertificateRequest) (*awsv1beta1.CertificateRequest, error) {
	result := awsv1beta1.CertificateRequest{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		Name(certificaterequest.Name).
		Body(certificaterequest).
		Do().
		Into(&result)

	return &result, err
}

// Delete a CertificateRequest.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		Name(name).
		Do().
		Error()
}

// Watch for CertificateRequests.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
