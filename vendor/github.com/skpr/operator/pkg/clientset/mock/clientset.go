package mock

import (
	"github.com/skpr/operator/pkg/clientset"
	"github.com/skpr/operator/pkg/clientset/app"
	mockapp "github.com/skpr/operator/pkg/clientset/app/mock"
	"github.com/skpr/operator/pkg/clientset/aws"
	mockaws "github.com/skpr/operator/pkg/clientset/aws/mock"
	"github.com/skpr/operator/pkg/clientset/edge"
	mockedge "github.com/skpr/operator/pkg/clientset/edge/mock"
	"github.com/skpr/operator/pkg/clientset/extensions"
	mockextensions "github.com/skpr/operator/pkg/clientset/extensions/mock"
	"github.com/skpr/operator/pkg/clientset/mysql"
	mockmysql "github.com/skpr/operator/pkg/clientset/mysql/mock"
	"k8s.io/apimachinery/pkg/runtime"
)

// New clientset using mock implementations.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	var (
		appObjects       []runtime.Object
		awsObjects       []runtime.Object
		edgeObjects      []runtime.Object
		mysqlObjects     []runtime.Object
		extensionObjects []runtime.Object
	)

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Group {
		case app.Group:
			appObjects = append(appObjects, object)
		case aws.Group:
			awsObjects = append(awsObjects, object)
		case edge.Group:
			edgeObjects = append(edgeObjects, object)
		case mysql.Group:
			mysqlObjects = append(mysqlObjects, object)
		case extensions.Group:
			extensionObjects = append(extensionObjects, object)
		}
	}

	appClient, err := mockapp.New(appObjects...)
	if err != nil {
		return nil, err
	}

	awsClient, err := mockaws.New(awsObjects...)
	if err != nil {
		return nil, err
	}

	edgeClient, err := mockedge.New(edgeObjects...)
	if err != nil {
		return nil, err
	}

	mysqlClient, err := mockmysql.New(mysqlObjects...)
	if err != nil {
		return nil, err
	}

	extensionsClient, err := mockextensions.New(extensionObjects...)
	if err != nil {
		return nil, err
	}

	clients := clientset.Clients{
		App:        appClient,
		AWS:        awsClient,
		Edge:       edgeClient,
		Mysql:      mysqlClient,
		Extensions: extensionsClient,
	}

	return clientset.New(clients)
}
