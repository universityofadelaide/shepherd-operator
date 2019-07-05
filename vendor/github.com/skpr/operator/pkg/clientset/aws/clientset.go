package aws

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// ClientInterface for interacting with AWS subclients.
type ClientInterface interface {
	Certificates(namespace string) CertificateInterface
	CertificateRequests(namespace string) CertificateRequestInterface
	CloudFronts(namespace string) CloudFrontInterface
	CloudFrontInvalidations(namespace string) CloudFrontInvalidationInterface
}

// Client for interacting with Operator objects.
type Client struct {
	RestClient rest.Interface
}

// NewForConfig returns a client for interacting with AWS objects.
func NewForConfig(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &awsv1beta1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{RestClient: client}, nil
}

// Certificates within a namespace.
func (c *Client) Certificates(namespace string) CertificateInterface {
	return &certificateClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CertificateRequests within a namespace.
func (c *Client) CertificateRequests(namespace string) CertificateRequestInterface {
	return &certificateRequestClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CloudFronts within a namespace.
func (c *Client) CloudFronts(namespace string) CloudFrontInterface {
	return &cloudfrontClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CloudFrontInvalidations within a namespace.
func (c *Client) CloudFrontInvalidations(namespace string) CloudFrontInvalidationInterface {
	return &cloudFrontInvalidationClient{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
