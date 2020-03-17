package backup

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Kind of object this client handles.
const Kind = "Backup"

// Interface declares interactions with the Backup objects.
type Interface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.BackupList, error)
	Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Backup, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error)
	Update(*extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Backup objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Backups.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.BackupList, error) {
	result := extensionsv1beta1.BackupList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a Backup.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Backup, error) {
	result := extensionsv1beta1.Backup{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a Backup is present.
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

// Create a Backup.
func (c *Client) Create(backup *extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error) {
	result := extensionsv1beta1.Backup{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("backups").
		Body(backup).
		Do().
		Into(&result)

	return &result, err
}

// Update a Backup.
func (c *Client) Update(backup *extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error) {
	result := extensionsv1beta1.Backup{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("backups").
		Name(backup.Name).
		Body(backup).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Backup.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("backups").
		Name(name).
		Do().
		Error()
}

// Watch for Backups.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
