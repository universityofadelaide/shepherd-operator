package database

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// Kind of object this client handles.
const Kind = "Database"

// Interface declares interactions with the Database objects.
type Interface interface {
	List(opts metav1.ListOptions) (*mysqlv1beta1.DatabaseList, error)
	Get(name string, opts metav1.GetOptions) (*mysqlv1beta1.Database, error)
	Exists(name string, opts metav1.GetOptions) (bool, error)
	Create(*mysqlv1beta1.Database) (*mysqlv1beta1.Database, error)
	Update(*mysqlv1beta1.Database) (*mysqlv1beta1.Database, error)
	Delete(name string, opts *metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

// Client for interacting with Database objects.
type Client struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Databasees.
func (c *Client) List(opts metav1.ListOptions) (*mysqlv1beta1.DatabaseList, error) {
	result := mysqlv1beta1.DatabaseList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databasees").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Get a Database.
func (c *Client) Get(name string, opts metav1.GetOptions) (*mysqlv1beta1.Database, error) {
	result := mysqlv1beta1.Database{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databasees").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
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
	result := mysqlv1beta1.Database{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("databasees").
		Body(database).
		Do().
		Into(&result)

	return &result, err
}

// Update a Database.
func (c *Client) Update(database *mysqlv1beta1.Database) (*mysqlv1beta1.Database, error) {
	result := mysqlv1beta1.Database{}
	err := c.RestClient.
		Put().
		Namespace(c.Namespace).
		Resource("databasees").
		Name(database.Name).
		Body(database).
		Do().
		Into(&result)

	return &result, err
}

// Delete a Database.
func (c *Client) Delete(name string, options *metav1.DeleteOptions) error {
	return c.RestClient.
		Delete().
		Namespace(c.Namespace).
		Resource("databasees").
		Name(name).
		Do().
		Error()
}

// Watch for Databasees.
func (c *Client) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databasees").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
