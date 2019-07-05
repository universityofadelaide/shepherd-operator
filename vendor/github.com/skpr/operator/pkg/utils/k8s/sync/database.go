package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// Database function to ensure the Object is in sync.
func Database(parent metav1.Object, spec mysqlv1beta1.DatabaseSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		database := obj.(*mysqlv1beta1.Database)
		database.Spec = spec
		return controllerutil.SetControllerReference(parent, database, scheme)
	}
}
