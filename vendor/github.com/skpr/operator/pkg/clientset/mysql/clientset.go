package mysql

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	"github.com/skpr/operator/pkg/clientset/mysql/database"
)

const (
	// Group which this clientset interacts with.
	Group = "mysql.skpr.io"
	// Version which this clientset interacts with.
	Version = "v1beta1"
	// APIVersion which this clientset interacts with.
	APIVersion = "mysql.skpr.io/v1beta1"
)

// Interface for interacting with AWS subclients.
type Interface interface {
	Databases(namespace string) database.Interface
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
func (c *Client) Databases(namespace string) database.Interface {
	return &database.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
