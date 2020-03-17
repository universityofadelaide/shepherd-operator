package sync

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	searchv1beta1 "github.com/skpr/operator/pkg/apis/search/v1beta1"
)

// Solr function to ensure the Object is in sync.
func Solr(parent metav1.Object, spec searchv1beta1.SolrSpec, scheme *runtime.Scheme) controllerutil.MutateFn {
	return func(obj runtime.Object) error {
		solr := obj.(*searchv1beta1.Solr)
		solr.Spec = spec
		return controllerutil.SetControllerReference(parent, solr, scheme)
	}
}
