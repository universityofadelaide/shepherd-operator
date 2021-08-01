package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
)

// Backup function to ensure the Object is in sync.
func Backup(parent metav1.Object, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		backup := obj.(*extensionv1.Backup)
		return controllerutil.SetControllerReference(parent, backup, scheme)
	}
}
