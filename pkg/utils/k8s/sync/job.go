package sync

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Job function to ensure the Object is in sync.
func Job(parent metav1.Object, spec batchv1.JobSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		job := obj.(*batchv1.Job)
		job.Spec = spec
		return controllerutil.SetControllerReference(parent, job, scheme)
	}
}
