package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
)

// Client for mocking Drupal client.
type Client struct {
	Namespace string
	Objects   []*appv1beta1.Drupal
}

// List a Drupal.
func (c *Client) List(opts metav1.ListOptions) (*appv1beta1.DrupalList, error) {
	result := &appv1beta1.DrupalList{}

	for _, drupal := range c.Objects {
		result.Items = append(result.Items, *drupal)
	}

	return result, nil
}

// Get a Drupal.
func (c *Client) Get(name string, opts metav1.GetOptions) (*appv1beta1.Drupal, error) {
	for _, drupal := range c.Objects {
		if drupal.ObjectMeta.Name == name {
			return drupal, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
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
	exists, err := c.Exists(drupal.ObjectMeta.Name, metav1.GetOptions{})
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

	drupal.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, drupal)

	return drupal, nil
}

// Update a Drupal.
func (c *Client) Update(drupal *appv1beta1.Drupal) (*appv1beta1.Drupal, error) {
	exists, err := c.Exists(drupal.ObjectMeta.Name, metav1.GetOptions{})
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

	drupal.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == drupal.ObjectMeta.Name {
			c.Objects[i] = drupal
		}
	}

	return drupal, nil
}

// Delete a Drupal.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*appv1beta1.Drupal

	for _, drupal := range c.Objects {
		if drupal.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, drupal)
	}

	c.Objects = filtered

	return nil
}

// Watch for Drupal object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
