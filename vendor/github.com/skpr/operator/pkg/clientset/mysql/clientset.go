package mysql

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// ClientInterface for interacting with AWS subclients.
type ClientInterface interface {
	Databases(namespace string) DatabaseInterface
}

// Client for interacting with Operator objects.
type Client struct {
	RestClient rest.Interface
}

// NewForConfig returns a client for interacting with MySQL objects.
func NewForConfig(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &mysqlv1beta1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{RestClient: client}, nil
}

// Databases within a namespace.
func (c *Client) Databases(namespace string) DatabaseInterface {
	return &databaseClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
