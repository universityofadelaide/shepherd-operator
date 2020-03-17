package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// Backup function to ensure the Object is in sync.
func Backup(parent metav1.Object, spec extensionsv1beta1.BackupSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		have := obj.(*extensionsv1beta1.Backup)
		have.Spec = spec
		return controllerutil.SetControllerReference(parent, have, scheme)
	}
}
