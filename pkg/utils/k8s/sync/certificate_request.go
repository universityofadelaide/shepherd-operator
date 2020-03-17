package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// CertificateRequest function to ensure the Object is in sync.
func CertificateRequest(parent metav1.Object, spec awsv1beta1.CertificateRequestSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		request := obj.(*awsv1beta1.CertificateRequest)
		// We don't want to change any of the spec. A new certificate should be created if the spec needs to change.
		return controllerutil.SetControllerReference(parent, request, scheme)
	}
}
