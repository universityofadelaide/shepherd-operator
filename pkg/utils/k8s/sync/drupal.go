package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
)

// Drupal function to ensure the Object is in sync.
func Drupal(parent metav1.Object, spec appv1beta1.DrupalSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		drupal := obj.(*appv1beta1.Drupal)
		drupal.Spec = spec
		return controllerutil.SetControllerReference(parent, drupal, scheme)
	}
}
