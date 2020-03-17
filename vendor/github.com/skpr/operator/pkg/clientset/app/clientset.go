package app

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
	"github.com/skpr/operator/pkg/clientset/app/drupal"
)

const (
	// Group which this clientset interacts with.
	Group = "apps.skpr.io"
	// Version which this clientset interacts with.
	Version = "v1beta1"
	// APIVersion which this clientset interacts with.
	APIVersion = "apps.skpr.io/v1beta1"
)

// Interface for interacting with AWS subclients.
type Interface interface {
	Drupals(namespace string) drupal.Interface
}

// Client for interacting with Operator objects.
type Client struct {
	RestClient rest.Interface
}

// NewForConfig returns a client for interacting with Operator objects.
func NewForConfig(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &appv1beta1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{RestClient: client}, nil
}

// Drupals within a namespace.
func (c *Client) Drupals(namespace string) drupal.Interface {
	return &drupal.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
