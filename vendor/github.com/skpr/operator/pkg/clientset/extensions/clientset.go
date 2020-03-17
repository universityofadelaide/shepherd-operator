package extensions

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	"github.com/skpr/operator/pkg/clientset/extensions/backup"
	"github.com/skpr/operator/pkg/clientset/extensions/backupscheduled"
	"github.com/skpr/operator/pkg/clientset/extensions/exec"
	"github.com/skpr/operator/pkg/clientset/extensions/restore"
)

const (
	// Group which this clientset interacts with.
	Group = "extensions.skpr.io"
	// Version which this clientset interacts with.
	Version = "v1beta1"
	// APIVersion which this clientset interacts with.
	APIVersion = "extensions.skpr.io/v1beta1"
)

// Interface for interacting with AWS subclients.
type Interface interface {
	Backups(namespace string) backup.Interface
	BackupScheduleds(namespace string) backupscheduled.Interface
	Restores(namespace string) restore.Interface
	Execs(namespace string) exec.Interface
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
func (c *Client) Backups(namespace string) backup.Interface {
	return &backup.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// Execs within a namespace.
func (c *Client) Execs(namespace string) exec.Interface {
	return &exec.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// BackupScheduleds within a namespace.
func (c *Client) BackupScheduleds(namespace string) backupscheduled.Interface {
	return &backupscheduled.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}

// Restores within a namespace.
func (c *Client) Restores(namespace string) restore.Interface {
	return &restore.Client{
		RestClient: c.RestClient,
		Namespace:  namespace,
	}
}
