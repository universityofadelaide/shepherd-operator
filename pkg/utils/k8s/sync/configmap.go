package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ConfigMap function to ensure the Object is in sync.
func ConfigMap(parent metav1.Object, data map[string]string, binary map[string][]byte, overwrite bool, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		configmap := obj.(*corev1.ConfigMap)
		if overwrite {
			configmap.Data = data
			configmap.BinaryData = binary
		}
		return controllerutil.SetControllerReference(parent, configmap, scheme)
	}
}
