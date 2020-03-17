package drupal

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
)

// Kind of object this client handles.
const Kind = "Drupal"

// Interface declares interactions with the Drupal objects.
type Interface interface {
	List(opts metav1.ListOptions) (*appv1beta1.DrupalList, error)
	Get(name string, opts metav1.GetOptions) (*appv1beta1.Drupal, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*appv1beta1.Drupal) (*appv1beta1.Drupal, error)
	Update(*appv1beta1.Drupal) (*appv1beta1.Drupal, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Drupal objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Drupals.
func (c *Client) List(opts metav1.ListOptions) (*appv1beta1.DrupalList, error) {
	result := appv1beta1.DrupalList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a Drupal.
func (c *Client) Get(name string, opts metav1.GetOptions) (*appv1beta1.Drupal, error) {
	result := appv1beta1.Drupal{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a Drupal is present.
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

// Create a Drupal.
func (c *Client) Create(drupal *appv1beta1.Drupal) (*appv1beta1.Drupal, error) {
	result := appv1beta1.Drupal{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("drupals").
		Body(drupal).
		Do().
		Into(&result)

	return &result, err
}

// Update a Drupal.
func (c *Client) Update(drupal *appv1beta1.Drupal) (*appv1beta1.Drupal, error) {
	result := appv1beta1.Drupal{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("drupals").
		Name(drupal.Name).
		Body(drupal).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Drupal.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("drupals").
		Name(name).
		Do().
		Error()
}

// Watch for Drupals.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
