package edge

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
)

// IngressInterface declares interactions with the Ingress objects.
type IngressInterface interface {
	List(opts metav1.ListOptions) (*edgev1beta1.IngressList, error)
	Get(name string, options metav1.GetOptions) (*edgev1beta1.Ingress, error)
	Create(*edgev1beta1.Ingress) (*edgev1beta1.Ingress, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type ingressClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Ingresss.
func (c *ingressClient) List(opts metav1.ListOptions) (*edgev1beta1.IngressList, error) {
	result := edgev1beta1.IngressList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("ingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Ingress.
func (c *ingressClient) Get(name string, opts metav1.GetOptions) (*edgev1beta1.Ingress, error) {
	result := edgev1beta1.Ingress{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("ingresses").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Ingress.
func (c *ingressClient) Create(ingress *edgev1beta1.Ingress) (*edgev1beta1.Ingress, error) {
	result := edgev1beta1.Ingress{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("ingresses").
		Body(ingress).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Ingresss.
func (c *ingressClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("ingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
