package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Secret function to ensure the Object is in sync.
func Secret(parent metav1.Object, data map[string][]byte, overwrite bool, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		secret := obj.(*corev1.Secret)
		if overwrite {
			secret.Data = data
		}
		return controllerutil.SetControllerReference(parent, secret, scheme)
	}
}
