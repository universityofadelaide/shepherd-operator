package mock

import (
	"k8s.io/apimachinery/pkg/runtime"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	clientset "github.com/skpr/operator/pkg/clientset/mysql"
	"github.com/skpr/operator/pkg/clientset/mysql/database"
	databasemock "github.com/skpr/operator/pkg/clientset/mysql/database/mock"
)

// Clientset used for mocking.
type Clientset struct {
	databases []*mysqlv1beta1.Database
}

// New clientset for interacting with Workflow objects.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	clientset := &Clientset{}

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Kind {
		case database.Kind:
			clientset.databases = append(clientset.databases, object.(*mysqlv1beta1.Database))
		}
	}

	return clientset, nil
}

// Databases within a namespace.
func (c *Clientset) Databases(namespace string) database.Interface {
	filter := func(list []*mysqlv1beta1.Database, namespace string) []*mysqlv1beta1.Database {
		var filtered []*mysqlv1beta1.Database

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &databasemock.Client{
		Namespace: namespace,
		Objects:   filter(c.databases, namespace),
	}
}
