package backupscheduled

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Kind of object this client handles.
const Kind = "BackupScheduled"

// Interface declares interactions with the BackupScheduled objects.
type Interface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.BackupScheduledList, error)
	Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.BackupScheduled, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error)
	Update(*extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with BackupScheduled objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all BackupScheduleds.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.BackupScheduledList, error) {
	result := extensionsv1beta1.BackupScheduledList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a BackupScheduled.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.BackupScheduled, error) {
	result := extensionsv1beta1.BackupScheduled{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Exists checks if a BackupScheduled is present.
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

// Create a BackupScheduled.
func (c *Client) Create(backupscheduled *extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error) {
	result := extensionsv1beta1.BackupScheduled{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		Body(backupscheduled).
		Do().
		Into(&result)

	return &result, err
}

// Update a BackupScheduled.
func (c *Client) Update(backupscheduled *extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error) {
	result := extensionsv1beta1.BackupScheduled{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		Name(backupscheduled.Name).
		Body(backupscheduled).
		Do().
		Into(&result)

	return &result, err
}

// Delete a BackupScheduled.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		Name(name).
		Do().
		Error()
}

// Watch for BackupScheduleds.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backupscheduleds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
