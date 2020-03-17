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
	App() app.Interface
	AWS() aws.Interface
	Edge() edge.Interface
	Extensions() extensions.Interface
	MySQL() mysql.Interface
}

// Clientset for interacting with Operator objects.
type Clientset struct {
	Clients
}

// Clients used as part of the clientset.
type Clients struct {
	App        app.Interface
	AWS        aws.Interface
	Edge       edge.Interface
	Extensions extensions.Interface
	Mysql      mysql.Interface
}

// New clientset using a set of clients.
func New(clients Clients) (Interface, error) {
	return &Clientset{clients}, nil
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

	clients := Clients{
		App:        appClient,
		AWS:        awsClient,
		Edge:       edgeClient,
		Extensions: extensionsClient,
		Mysql:      mysqlClient,
	}

	return New(clients)
}

// App clientset.
func (c *Clientset) App() app.Interface {
	return c.Clients.App
}

// AWS clientset.
func (c *Clientset) AWS() aws.Interface {
	return c.Clients.AWS
}

// Edge clientset.
func (c *Clientset) Edge() edge.Interface {
	return c.Clients.Edge
}

// Extensions clientset.
func (c *Clientset) Extensions() extensions.Interface {
	return c.Clients.Extensions
}

// MySQL clientset.
func (c *Clientset) MySQL() mysql.Interface {
	return c.Clients.Mysql
}

func init() {
	apis.AddToScheme(scheme.Scheme)
}
