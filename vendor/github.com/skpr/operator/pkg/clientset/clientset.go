package clientset

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/skpr/operator/pkg/apis"
	"github.com/skpr/operator/pkg/clientset/app"
	"github.com/skpr/operator/pkg/clientset/aws"
	"github.com/skpr/operator/pkg/clientset/edge"
	"github.com/skpr/operator/pkg/clientset/extensions"
	"github.com/skpr/operator/pkg/clientset/mysql"
)

// Interface for interacting with Operator subclients.
type Interface interface {
	App() app.ClientInterface
	AWS() aws.ClientInterface
	Edge() edge.ClientInterface
	Extensions() extensions.ClientInterface
	MySQL() mysql.ClientInterface
}

// Client for interacting with Operator objects.
type Client struct {
	app        app.ClientInterface
	aws        aws.ClientInterface
	edge       edge.ClientInterface
	extensions extensions.ClientInterface
	mysql      mysql.ClientInterface
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(config *rest.Config) (Interface, error) {
	appClient, err := app.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	awsClient, err := aws.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	edgeClient, err := edge.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	extensionsClient, err := extensions.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	mysqlClient, err := mysql.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		app:        appClient,
		aws:        awsClient,
		edge:       edgeClient,
		extensions: extensionsClient,
		mysql:      mysqlClient,
	}, nil
}

// App clientset.
func (c *Client) App() app.ClientInterface {
	return c.app
}

// AWS clientset.
func (c *Client) AWS() aws.ClientInterface {
	return c.aws
}

// Edge clientset.
func (c *Client) Edge() edge.ClientInterface {
	return c.edge
}

// Extensions clientset.
func (c *Client) Extensions() extensions.ClientInterface {
	return c.extensions
}

// MySQL clientset.
func (c *Client) MySQL() mysql.ClientInterface {
	return c.mysql
}

func init() {
	apis.AddToScheme(scheme.Scheme)
}
