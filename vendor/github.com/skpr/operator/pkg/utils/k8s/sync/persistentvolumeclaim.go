package sync

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// PersistentVolumeClaim function to ensure the Object is in sync.
func PersistentVolumeClaim(parent metav1.Object, spec corev1.PersistentVolumeClaimSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		pvc := obj.(*corev1.PersistentVolumeClaim)
		return controllerutil.SetControllerReference(parent, pvc, scheme)
	}
}
