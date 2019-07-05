package extensions

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// ClientInterface for interacting with AWS subclients.
type ClientInterface interface {
	Backups(namespace string) BackupInterface
	Execs(namespace string) ExecInterface
}

// Client for interacting with Operator objects.
type Client struct {
	RestClient rest.Interface
}

// NewForConfig returns a client for interacting with Extensions objects.
func NewForConfig(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &extensionsv1beta1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{RestClient: client}, nil
}

// Backups within a namespace.
func (c *Client) Backups(namespace string) BackupInterface {
	return &backupClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// Execs within a namespace.
func (c *Client) Execs(namespace string) ExecInterface {
	return &execClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
