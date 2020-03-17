package mock

import (
	"k8s.io/apimachinery/pkg/runtime"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
	clientset "github.com/skpr/operator/pkg/clientset/app"
	"github.com/skpr/operator/pkg/clientset/app/drupal"
	drupalmock "github.com/skpr/operator/pkg/clientset/app/drupal/mock"
)

// Clientset used for mocking.
type Clientset struct {
	drupals []*appv1beta1.Drupal
}

// New clientset for interacting with Workflow objects.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	clientset := &Clientset{}

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Kind {
		case drupal.Kind:
			clientset.drupals = append(clientset.drupals, object.(*appv1beta1.Drupal))
		}
	}

	return clientset, nil
}

// Drupals within a namespace.
func (c *Clientset) Drupals(namespace string) drupal.Interface {
	filter := func(list []*appv1beta1.Drupal, namespace string) []*appv1beta1.Drupal {
		var filtered []*appv1beta1.Drupal

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &drupalmock.Client{
		Namespace: namespace,
		Objects:   filter(c.drupals, namespace),
	}
}
