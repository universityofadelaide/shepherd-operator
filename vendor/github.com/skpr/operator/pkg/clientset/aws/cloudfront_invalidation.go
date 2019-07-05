package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// CloudFrontInvalidationInterface declares interactions with the CloudFrontInvalidation objects.
type CloudFrontInvalidationInterface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontInvalidationList, error)
	Get(name string, options metav1.GetOptions) (*awsv1beta1.CloudFrontInvalidation, error)
	Create(*awsv1beta1.CloudFrontInvalidation) (*awsv1beta1.CloudFrontInvalidation, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type cloudFrontInvalidationClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all CloudFrontInvalidations.
func (c *cloudFrontInvalidationClient) List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontInvalidationList, error) {
	result := awsv1beta1.CloudFrontInvalidationList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an CloudFrontInvalidation.
func (c *cloudFrontInvalidationClient) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CloudFrontInvalidation, error) {
	result := awsv1beta1.CloudFrontInvalidation{}
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

// Create an CloudFrontInvalidation.
func (c *cloudFrontInvalidationClient) Create(certificaterequest *awsv1beta1.CloudFrontInvalidation) (*awsv1beta1.CloudFrontInvalidation, error) {
	result := awsv1beta1.CloudFrontInvalidation{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		Body(certificaterequest).
		Do().
		Into(&result)

	return &result, err
}

// Watch for CloudFrontInvalidations.
func (c *cloudFrontInvalidationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
