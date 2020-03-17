package sync

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// StatefulSet function to ensure the Object is in sync.
func StatefulSet(parent metav1.Object, spec appsv1.StatefulSetSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		sts := obj.(*appsv1.StatefulSet)
		sts.Spec.Selector = spec.Selector
		sts.Spec.Template = spec.Template
		return controllerutil.SetControllerReference(parent, sts, scheme)
	}
}
