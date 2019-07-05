package mysql

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// DatabaseInterface declares interactions with the Database objects.
type DatabaseInterface interface {
	List(opts metav1.ListOptions) (*mysqlv1beta1.DatabaseList, error)
	Get(name string, options metav1.GetOptions) (*mysqlv1beta1.Database, error)
	Create(*mysqlv1beta1.Database) (*mysqlv1beta1.Database, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type databaseClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Databases.
func (c *databaseClient) List(opts metav1.ListOptions) (*mysqlv1beta1.DatabaseList, error) {
	result := mysqlv1beta1.DatabaseList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Database.
func (c *databaseClient) Get(name string, opts metav1.GetOptions) (*mysqlv1beta1.Database, error) {
	result := mysqlv1beta1.Database{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databases").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Database.
func (c *databaseClient) Create(database *mysqlv1beta1.Database) (*mysqlv1beta1.Database, error) {
	result := mysqlv1beta1.Database{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("databases").
		Body(database).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Databases.
func (c *databaseClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("databases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
