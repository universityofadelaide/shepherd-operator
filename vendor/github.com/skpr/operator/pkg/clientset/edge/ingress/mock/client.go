package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
)

// Client for mocking Ingress client.
type Client struct {
	Namespace string
	Objects   []*edgev1beta1.Ingress
}

// List a Ingress.
func (c *Client) List(opts metav1.ListOptions) (*edgev1beta1.IngressList, error) {
	result := &edgev1beta1.IngressList{}

	for _, ingress := range c.Objects {
		result.Items = append(result.Items, *ingress)
	}

	return result, nil
}

// Get a Ingress.
func (c *Client) Get(name string, opts metav1.GetOptions) (*edgev1beta1.Ingress, error) {
	for _, ingress := range c.Objects {
		if ingress.ObjectMeta.Name == name {
			return ingress, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
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
	exists, err := c.Exists(ingress.ObjectMeta.Name, metav1.GetOptions{})
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

	ingress.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, ingress)

	return ingress, nil
}

// Update a Ingress.
func (c *Client) Update(ingress *edgev1beta1.Ingress) (*edgev1beta1.Ingress, error) {
	exists, err := c.Exists(ingress.ObjectMeta.Name, metav1.GetOptions{})
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

	ingress.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == ingress.ObjectMeta.Name {
			c.Objects[i] = ingress
		}
	}

	return ingress, nil
}

// Delete a Ingress.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*edgev1beta1.Ingress

	for _, ingress := range c.Objects {
		if ingress.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, ingress)
	}

	c.Objects = filtered

	return nil
}

// Watch for Ingress object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
