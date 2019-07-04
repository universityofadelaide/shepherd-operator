package sync

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Job function to ensure the Object is in sync.
func Job(parent metav1.Object, job *batchv1.Job, scheme *runtime.Scheme) controllerutil.MutateFn {
	err := func() error {
		return controllerutil.SetControllerReference(parent, job, scheme)
	}

	// @todo fix this error handling and return value.
	if err != nil {
		panic(err)
	}
	return nil
}
