package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Client for mocking CloudFrontInvalidation client.
type Client struct {
	Namespace string
	Objects   []*awsv1beta1.CloudFrontInvalidation
}

// List a CloudFrontInvalidation.
func (c *Client) List(opts metav1.ListOptions) (*awsv1beta1.CloudFrontInvalidationList, error) {
	result := &awsv1beta1.CloudFrontInvalidationList{}

	for _, cloudFront := range c.Objects {
		result.Items = append(result.Items, *cloudFront)
	}

	return result, nil
}

// Get a CloudFrontInvalidation.
func (c *Client) Get(name string, opts metav1.GetOptions) (*awsv1beta1.CloudFrontInvalidation, error) {
	for _, cloudFront := range c.Objects {
		if cloudFront.ObjectMeta.Name == name {
			return cloudFront, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
}

// Exists checks if a CloudFrontInvalidation is present.
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

// Create a CloudFrontInvalidation.
func (c *Client) Create(cloudFront *awsv1beta1.CloudFrontInvalidation) (*awsv1beta1.CloudFrontInvalidation, error) {
	exists, err := c.Exists(cloudFront.ObjectMeta.Name, metav1.GetOptions{})
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

	cloudFront.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, cloudFront)

	return cloudFront, nil
}

// Update a CloudFrontInvalidation.
func (c *Client) Update(cloudFront *awsv1beta1.CloudFrontInvalidation) (*awsv1beta1.CloudFrontInvalidation, error) {
	exists, err := c.Exists(cloudFront.ObjectMeta.Name, metav1.GetOptions{})
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

	cloudFront.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == cloudFront.ObjectMeta.Name {
			c.Objects[i] = cloudFront
		}
	}

	return cloudFront, nil
}

// Delete a CloudFrontInvalidation.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*awsv1beta1.CloudFrontInvalidation

	for _, cloudFront := range c.Objects {
		if cloudFront.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, cloudFront)
	}

	c.Objects = filtered

	return nil
}

// Watch for CloudFrontInvalidation object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
