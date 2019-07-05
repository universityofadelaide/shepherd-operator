package edge

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
)

// ClientInterface for interacting with AWS subclients.
type ClientInterface interface {
	Ingresss(namespace string) IngressInterface
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

// Ingresss within a namespace.
func (c *Client) Ingresss(namespace string) IngressInterface {
	return &ingressClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
