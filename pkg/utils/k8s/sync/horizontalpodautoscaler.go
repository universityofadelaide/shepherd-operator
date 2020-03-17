package sync

import (
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// HPA function to ensure the Object is in sync.
func HPA(parent metav1.Object, spec autoscalingv2beta2.HorizontalPodAutoscalerSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		hpa := obj.(*autoscalingv2beta2.HorizontalPodAutoscaler)
		hpa.Spec = spec
		return controllerutil.SetControllerReference(parent, hpa, scheme)
	}
}
