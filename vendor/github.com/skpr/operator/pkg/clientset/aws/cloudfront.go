package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// CloudFrontInterface declares interactions with the CloudFront objects.
type CloudFrontInterface interface {
	List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontList, error)
	Get(name string, options metav1.GetOptions) (*awsv1beta1.CloudFront, error)
	Create(*awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type cloudfrontClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all CloudFronts.
func (c *cloudfrontClient) List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontList, error) {
	result := awsv1beta1.CloudFrontList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudfronts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an CloudFront.
func (c *cloudfrontClient) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CloudFront, error) {
	result := awsv1beta1.CloudFront{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudfronts").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an CloudFront.
func (c *cloudfrontClient) Create(cloudfront *awsv1beta1.CloudFront) (*awsv1beta1.CloudFront, error) {
	result := awsv1beta1.CloudFront{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("cloudfronts").
		Body(cloudfront).
		Do().
		Into(&result)

	return &result, err
}

// Watch for CloudFronts.
func (c *cloudfrontClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("cloudfronts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
