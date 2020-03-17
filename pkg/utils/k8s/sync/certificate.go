package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Certificate function to ensure the Object is in sync.
func Certificate(parent metav1.Object, spec awsv1beta1.CertificateSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		ingress := obj.(*awsv1beta1.Certificate)
		ingress.Spec = spec
		return controllerutil.SetControllerReference(parent, ingress, scheme)
	}
}
