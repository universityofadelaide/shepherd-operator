package extensions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// ExecInterface declares interactions with the Exec objects.
type ExecInterface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.ExecList, error)
	Get(name string, options metav1.GetOptions) (*extensionsv1beta1.Exec, error)
	Create(*extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type execClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Execs.
func (c *execClient) List(opts metav1.ListOptions) (*extensionsv1beta1.ExecList, error) {
	result := extensionsv1beta1.ExecList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Exec.
func (c *execClient) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Exec, error) {
	result := extensionsv1beta1.Exec{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Exec.
func (c *execClient) Create(exec *extensionsv1beta1.Exec) (*extensionsv1beta1.Exec, error) {
	result := extensionsv1beta1.Exec{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("execs").
		Body(exec).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Execs.
func (c *execClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("execs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
