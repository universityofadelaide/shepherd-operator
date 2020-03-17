package restore

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Kind of object this client handles.
const Kind = "Restore"

// Interface declares interactions with the Restore objects.
type Interface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.RestoreList, error)
	Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Restore, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*extensionsv1beta1.Restore) (*extensionsv1beta1.Restore, error)
	Update(*extensionsv1beta1.Restore) (*extensionsv1beta1.Restore, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Restore objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Restores.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.RestoreList, error) {
	result := extensionsv1beta1.RestoreList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("restores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a Restore.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Restore, error) {
	result := extensionsv1beta1.Restore{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("restores").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a Restore is present.
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

// Create a Restore.
func (c *Client) Create(restore *extensionsv1beta1.Restore) (*extensionsv1beta1.Restore, error) {
	result := extensionsv1beta1.Restore{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("restores").
		Body(restore).
		Do().
		Into(&result)

	return &result, err
}

// Update a Restore.
func (c *Client) Update(restore *extensionsv1beta1.Restore) (*extensionsv1beta1.Restore, error) {
	result := extensionsv1beta1.Restore{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("restores").
		Name(restore.Name).
		Body(restore).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Restore.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("restores").
		Name(name).
		Do().
		Error()
}

// Watch for Restores.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("restores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
