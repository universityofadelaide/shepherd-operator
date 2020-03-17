package ingress

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
)

// Kind of object this client handles.
const Kind = "Ingress"

// Interface declares interactions with the Ingress objects.
type Interface interface {
	List(opts metav1.ListOptions) (*edgev1beta1.IngressList, error)
	Get(name string, opts metav1.GetOptions) (*edgev1beta1.Ingress, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*edgev1beta1.Ingress) (*edgev1beta1.Ingress, error)
	Update(*edgev1beta1.Ingress) (*edgev1beta1.Ingress, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Ingress objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Ingresses.
func (c *Client) List(opts metav1.ListOptions) (*edgev1beta1.IngressList, error) {
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

// Get a Ingress.
func (c *Client) Get(name string, opts metav1.GetOptions) (*edgev1beta1.Ingress, error) {
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

// Exists checks if a Ingress is present.
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

// Create a Ingress.
func (c *Client) Create(ingress *edgev1beta1.Ingress) (*edgev1beta1.Ingress, error) {
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

// Update a Ingress.
func (c *Client) Update(ingress *edgev1beta1.Ingress) (*edgev1beta1.Ingress, error) {
	result := edgev1beta1.Ingress{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("ingresses").
		Name(ingress.Name).
		Body(ingress).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Ingress.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("ingresses").
		Name(name).
		Do().
		Error()
}

// Watch for Ingresses.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("ingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
