package sync

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Deployment function to ensure the Object is in sync.
func Deployment(parent metav1.Object, spec appsv1.DeploymentSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		deployment := obj.(*appsv1.Deployment)
		deployment.Spec.Selector = spec.Selector
		deployment.Spec.Template = spec.Template
		return controllerutil.SetControllerReference(parent, deployment, scheme)
	}
}
