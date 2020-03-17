package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Client for mocking BackupScheduled client.
type Client struct {
	Namespace string
	Objects   []*extensionsv1beta1.BackupScheduled
}

// List a BackupScheduled.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.BackupScheduledList, error) {
	result := &extensionsv1beta1.BackupScheduledList{}

	for _, backup := range c.Objects {
		result.Items = append(result.Items, *backup)
	}

	return result, nil
}

// Get a BackupScheduled.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.BackupScheduled, error) {
	for _, backup := range c.Objects {
		if backup.ObjectMeta.Name == name {
			return backup, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
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
func (c *Client) Create(backup *extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error) {
	exists, err := c.Exists(backup.ObjectMeta.Name, metav1.GetOptions{})
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

	backup.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, backup)

	return backup, nil
}

// Update a BackupScheduled.
func (c *Client) Update(backup *extensionsv1beta1.BackupScheduled) (*extensionsv1beta1.BackupScheduled, error) {
	exists, err := c.Exists(backup.ObjectMeta.Name, metav1.GetOptions{})
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

	backup.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == backup.ObjectMeta.Name {
			c.Objects[i] = backup
		}
	}

	return backup, nil
}

// Delete a BackupScheduled.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*extensionsv1beta1.BackupScheduled

	for _, backup := range c.Objects {
		if backup.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, backup)
	}

	c.Objects = filtered

	return nil
}

// Watch for BackupScheduled object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
