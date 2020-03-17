package edge

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
	"github.com/skpr/operator/pkg/clientset/edge/ingress"
)

const (
	// Group which this clientset interacts with.
	Group = "edge.skpr.io"
	// Version which this clientset interacts with.
	Version = "v1beta1"
	// APIVersion which this clientset interacts with.
	APIVersion = "edge.skpr.io/v1beta1"
)

// Interface for interacting with AWS subclients.
type Interface interface {
	Ingresses(namespace string) ingress.Interface
}

// Client for interacting with Operator objects.
type Client struct {
	RestClient rest.Interface
}

// NewForConfig returns a client for interacting with Edge objects.
func NewForConfig(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &edgev1beta1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{RestClient: client}, nil
}

// Ingresses within a namespace.
func (c *Client) Ingresses(namespace string) ingress.Interface {
	return &ingress.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
