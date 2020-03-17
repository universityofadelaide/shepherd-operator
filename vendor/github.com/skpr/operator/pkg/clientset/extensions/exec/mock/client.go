package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Client for mocking Exec client.
type Client struct {
	Namespace string
	Objects   []*extensionsv1beta1.Exec
}

// List a Exec.
func (c *Client) List(opts metav1.ListOptions) (*extensionsv1beta1.ExecList, error) {
	result := &extensionsv1beta1.ExecList{}

	for _, exec := range c.Objects {
		result.Items = append(result.Items, *exec)
	}

	return result, nil
}

// Get a Exec.
func (c *Client) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Exec, error) {
	for _, exec := range c.Objects {
		if exec.ObjectMeta.Name == name {
			return exec, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
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
	exists, err := c.Exists(exec.ObjectMeta.Name, metav1.GetOptions{})
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

	exec.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, exec)

	return exec, nil
}

// Update a Exec.
func (c *Client) Update(exec *extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error) {
	exists, err := c.Exists(exec.ObjectMeta.Name, metav1.GetOptions{})
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

	exec.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == exec.ObjectMeta.Name {
			c.Objects[i] = exec
		}
	}

	return exec, nil
}

// Delete a Exec.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*extensionsv1beta1.Exec

	for _, exec := range c.Objects {
		if exec.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, exec)
	}

	c.Objects = filtered

	return nil
}

// Watch for Exec object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
