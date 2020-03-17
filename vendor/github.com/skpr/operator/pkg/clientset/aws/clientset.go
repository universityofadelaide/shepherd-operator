package aws

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	"github.com/skpr/operator/pkg/clientset/aws/certificate"
	"github.com/skpr/operator/pkg/clientset/aws/certificaterequest"
	"github.com/skpr/operator/pkg/clientset/aws/cloudfront"
	"github.com/skpr/operator/pkg/clientset/aws/cloudfrontinvalidation"
)

const (
	// Group which this clientset interacts with.
	Group = "aws.skpr.io"
	// Version which this clientset interacts with.
	Version = "v1beta1"
	// APIVersion which this clientset interacts with.
	APIVersion = "aws.skpr.io/v1beta1"
)

// Interface for interacting with AWS subclients.
type Interface interface {
	Certificates(namespace string) certificate.Interface
	CertificateRequests(namespace string) certificaterequest.Interface
	CloudFronts(namespace string) cloudfront.Interface
	CloudFrontInvalidations(namespace string) cloudfrontinvalidation.Interface
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
func (c *Client) Certificates(namespace string) certificate.Interface {
	return &certificate.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CertificateRequests within a namespace.
func (c *Client) CertificateRequests(namespace string) certificaterequest.Interface {
	return &certificaterequest.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CloudFronts within a namespace.
func (c *Client) CloudFronts(namespace string) cloudfront.Interface {
	return &cloudfront.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// CloudFrontInvalidations within a namespace.
func (c *Client) CloudFrontInvalidations(namespace string) cloudfrontinvalidation.Interface {
	return &cloudfrontinvalidation.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
