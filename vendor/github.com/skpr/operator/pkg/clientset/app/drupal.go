package app

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
)

// DrupalInterface declares interactions with the Drupal objects.
type DrupalInterface interface {
	List(opts metav1.ListOptions) (*appv1beta1.DrupalList, error)
	Get(name string, options metav1.GetOptions) (*appv1beta1.Drupal, error)
	Create(*appv1beta1.Drupal) (*appv1beta1.Drupal, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type drupalClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Drupals.
func (c *drupalClient) List(opts metav1.ListOptions) (*appv1beta1.DrupalList, error) {
	result := appv1beta1.DrupalList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Drupal.
func (c *drupalClient) Get(name string, opts metav1.GetOptions) (*appv1beta1.Drupal, error) {
	result := appv1beta1.Drupal{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Drupal.
func (c *drupalClient) Create(drupal *appv1beta1.Drupal) (*appv1beta1.Drupal, error) {
	result := appv1beta1.Drupal{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("drupals").
		Body(drupal).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Drupals.
func (c *drupalClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("drupals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
