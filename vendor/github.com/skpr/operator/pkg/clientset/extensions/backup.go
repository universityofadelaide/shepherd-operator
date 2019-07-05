package extensions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// BackupInterface declares interactions with the Backup objects.
type BackupInterface interface {
	List(opts metav1.ListOptions) (*extensionsv1beta1.BackupList, error)
	Get(name string, options metav1.GetOptions) (*extensionsv1beta1.Backup, error)
	Create(*extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type backupClient struct {
	RestClient rest.Interface
	Namespace  string
}

// List all Backups.
func (c *backupClient) List(opts metav1.ListOptions) (*extensionsv1beta1.BackupList, error) {
	result := extensionsv1beta1.BackupList{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// List an Backup.
func (c *backupClient) Get(name string, opts metav1.GetOptions) (*extensionsv1beta1.Backup, error) {
	result := extensionsv1beta1.Backup{}
	err := c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

// Create an Backup.
func (c *backupClient) Create(backup *extensionsv1beta1.Backup) (*extensionsv1beta1.Backup, error) {
	result := extensionsv1beta1.Backup{}
	err := c.RestClient.
		Post().
		Namespace(c.Namespace).
		Resource("backups").
		Body(backup).
		Do().
		Into(&result)

	return &result, err
}

// Watch for Backups.
func (c *backupClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.RestClient.
		Get().
		Namespace(c.Namespace).
		Resource("backups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
