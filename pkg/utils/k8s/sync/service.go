package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Service function to ensure the Object is in sync.
func Service(parent metav1.Object, spec corev1.ServiceSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		service := obj.(*corev1.Service)
		service.Spec.Type = spec.Type
		service.Spec.Ports = spec.Ports
		service.Spec.Selector = spec.Selector
		return controllerutil.SetControllerReference(parent, service, scheme)
	}
}
