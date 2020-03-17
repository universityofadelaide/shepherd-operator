package cloudfront

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Kind of object this client handles.
const Kind = "CloudFront"

// Interface declares interactions with the CloudFront objects.
type Interface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontList, error)
	Get(name string, opts metav1.GetOptions) (*awsv1beta1.CloudFront, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error)
	Update(*awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with CloudFront objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all CloudFronts.
func (c *Client) List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontList, error) {
	result := awsv1beta1.CloudFrontList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a CloudFront.
func (c *Client) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CloudFront, error) {
	result := awsv1beta1.CloudFront{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a CloudFront is present.
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

// Create a CloudFront.
func (c *Client) Create(cloudFront *awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error) {
	result := awsv1beta1.CloudFront{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		Body(cloudFront).
		Do().
		Into(&result)

	return &result, err
}

// Update a CloudFront.
func (c *Client) Update(cloudFront *awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error) {
	result := awsv1beta1.CloudFront{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		Name(cloudFront.Name).
		Body(cloudFront).
		Do().
		Into(&result)

	return &result, err
}

// Delete a CloudFront.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		Name(name).
		Do().
		Error()
}

// Watch for CloudFronts.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudFronts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
