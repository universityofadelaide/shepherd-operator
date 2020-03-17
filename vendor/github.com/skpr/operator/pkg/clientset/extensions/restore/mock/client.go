package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Client for mocking Restore client.
type Client struct {
	Namespace string
	Objects   []*extensionsv1beta1.Restore
}

// List a Restore.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.RestoreList, error) {
	result := &extensionsv1beta1.RestoreList{}

	for _, restore := range c.Objects {
		result.Items = append(result.Items, *restore)
	}

	return result, nil
}

// Get a Restore.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Restore, error) {
	for _, restore := range c.Objects {
		if restore.ObjectMeta.Name == name {
			return restore, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
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
	exists, err := c.Exists(restore.ObjectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, &kerrors.StatusError{
			ErrStatus: metav1.Status{
				Reason: metav1.StatusReasonAlreadyExists,
			},
		}
	}

	restore.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, restore)

	return restore, nil
}

// Update a Restore.
func (c *Client) Update(restore *extensionsv1beta1.Restore) (*extensionsv1beta1.Restore, error) {
	exists, err := c.Exists(restore.ObjectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, &kerrors.StatusError{
			ErrStatus: metav1.Status{
				Reason: metav1.StatusReasonNotFound,
			},
		}
	}

	restore.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == restore.ObjectMeta.Name {
			c.Objects[i] = restore
		}
	}

	return restore, nil
}

// Delete a Restore.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*extensionsv1beta1.Restore

	for _, restore := range c.Objects {
		if restore.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, restore)
	}

	c.Objects = filtered

	return nil
}

// Watch for Restore object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
