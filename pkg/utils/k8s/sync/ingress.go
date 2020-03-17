package sync

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Ingress function to ensure the Object is in sync.
func Ingress(parent metav1.Object, want extensionsv1beta1.Ingress, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		ingress := obj.(*extensionsv1beta1.Ingress)
		ingress.ObjectMeta.Annotations = want.ObjectMeta.Annotations
		ingress.Spec = want.Spec
		return controllerutil.SetControllerReference(parent, ingress, scheme)
	}
}
