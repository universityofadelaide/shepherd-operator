package exec

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Kind of object this client handles.
const Kind = "Exec"

// Interface declares interactions with the Exec objects.
type Interface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.ExecList, error)
	Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Exec, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error)
	Update(*extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Exec objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Execs.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.ExecList, error) {
	result := extensionsv1beta1.ExecList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a Exec.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Exec, error) {
	result := extensionsv1beta1.Exec{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a Exec is present.
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

// Create a Exec.
func (c *Client) Create(exec *extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error) {
	result := extensionsv1beta1.Exec{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("execs").
		Body(exec).
		Do().
		Into(&result)

	return &result, err
}

// Update a Exec.
func (c *Client) Update(exec *extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error) {
	result := extensionsv1beta1.Exec{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("execs").
		Name(exec.Name).
		Body(exec).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Exec.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("execs").
		Name(name).
		Do().
		Error()
}

// Watch for Execs.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}