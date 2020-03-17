package mock

import (
	"k8s.io/apimachinery/pkg/runtime"

	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
	clientset "github.com/skpr/operator/pkg/clientset/edge"
	"github.com/skpr/operator/pkg/clientset/edge/ingress"
	ingressmock "github.com/skpr/operator/pkg/clientset/edge/ingress/mock"
)

// Clientset used for mocking.
type Clientset struct {
	ingresses []*edgev1beta1.Ingress
}

// New clientset for interacting with Workflow objects.
func New(objects ...runtime.Object) (clientset.Interface, error) {
	clientset := &Clientset{}

	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()

		switch gvk.Kind {
		case ingress.Kind:
			clientset.ingresses = append(clientset.ingresses, object.(*edgev1beta1.Ingress))
		}
	}

	return clientset, nil
}

// Ingresses within a namespace.
func (c *Clientset) Ingresses(namespace string) ingress.Interface {
	filter := func(list []*edgev1beta1.Ingress, namespace string) []*edgev1beta1.Ingress {
		var filtered []*edgev1beta1.Ingress

		for _, item := range list {
			if item.ObjectMeta.Namespace == namespace {
				filtered = append(filtered, item)
			}
		}

		return filtered
	}

	return &ingressmock.Client{
		Namespace: namespace,
		Objects:   filter(c.ingresses, namespace),
	}
}
