package sync

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CronJob function to ensure the Object is in sync.
func CronJob(parent metav1.Object, spec batchv1beta1.CronJobSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		cronjob := obj.(*batchv1beta1.CronJob)
		cronjob.Spec = spec
		return controllerutil.SetControllerReference(parent, cronjob, scheme)
	}
}
