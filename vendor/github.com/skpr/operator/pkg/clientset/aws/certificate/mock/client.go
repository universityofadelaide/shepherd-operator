package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Client for mocking Certificate client.
type Client struct {
	Namespace string
	Objects   []*awsv1beta1.Certificate
}

// List a Certificate.
func (c *Client) List(opts metav1.ListOptions) (*awsv1beta1.CertificateList, error) {
	result := &awsv1beta1.CertificateList{}

	for _, certificate := range c.Objects {
		result.Items = append(result.Items, *certificate)
	}

	return result, nil
}

// Get a Certificate.
func (c *Client) Get(name string, opts metav1.GetOptions) (*awsv1beta1.Certificate, error) {
	for _, certificate := range c.Objects {
		if certificate.ObjectMeta.Name == name {
			return certificate, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
}

// Exists checks if a Certificate is present.
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

// Create a Certificate.
func (c *Client) Create(certificate *awsv1beta1.Certificate) (*awsv1beta1.Certificate, error) {
	exists, err := c.Exists(certificate.ObjectMeta.Name, metav1.GetOptions{})
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

	certificate.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, certificate)

	return certificate, nil
}

// Update a Certificate.
func (c *Client) Update(certificate *awsv1beta1.Certificate) (*awsv1beta1.Certificate, error) {
	exists, err := c.Exists(certificate.ObjectMeta.Name, metav1.GetOptions{})
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

	certificate.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == certificate.ObjectMeta.Name {
			c.Objects[i] = certificate
		}
	}

	return certificate, nil
}

// Delete a Certificate.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*awsv1beta1.Certificate

	for _, certificate := range c.Objects {
		if certificate.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, certificate)
	}

	c.Objects = filtered

	return nil
}

// Watch for Certificate object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
