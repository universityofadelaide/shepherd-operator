package mock

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// Client for mocking Database client.
type Client struct {
	Namespace string
	Objects   []*mysqlv1beta1.Database
}

// List a Database.
func (c *Client) List(opts metav1.ListOptions) (*mysqlv1beta1.DatabaseList, error) {
	result := &mysqlv1beta1.DatabaseList{}

	for _, database := range c.Objects {
		result.Items = append(result.Items, *database)
	}

	return result, nil
}

// Get a Database.
func (c *Client) Get(name string, opts metav1.GetOptions) (*mysqlv1beta1.Database, error) {
	for _, database := range c.Objects {
		if database.ObjectMeta.Name == name {
			return database, nil
		}
	}

	return nil, &kerrors.StatusError{
		ErrStatus: metav1.Status{
			Reason: metav1.StatusReasonNotFound,
		},
	}
}

// Exists checks if a Database is present.
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

// Create a Database.
func (c *Client) Create(database *mysqlv1beta1.Database) (*mysqlv1beta1.Database, error) {
	exists, err := c.Exists(database.ObjectMeta.Name, metav1.GetOptions{})
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

	database.ObjectMeta.Namespace = c.Namespace

	c.Objects = append(c.Objects, database)

	return database, nil
}

// Update a Database.
func (c *Client) Update(database *mysqlv1beta1.Database) (*mysqlv1beta1.Database, error) {
	exists, err := c.Exists(database.ObjectMeta.Name, metav1.GetOptions{})
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

	database.ObjectMeta.Namespace = c.Namespace

	for i, current := range c.Objects {
		if current.ObjectMeta.Name == database.ObjectMeta.Name {
			c.Objects[i] = database
		}
	}

	return database, nil
}

// Delete a Database.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	var filtered []*mysqlv1beta1.Database

	for _, database := range c.Objects {
		if database.ObjectMeta.Name == name {
			continue
		}

		filtered = append(filtered, database)
	}

	c.Objects = filtered

	return nil
}

// Watch for Database object changes.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
