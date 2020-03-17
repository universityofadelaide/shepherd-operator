// +build unit

package sync

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/skpr/operator/pkg/apis"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	searchv1beta1 "github.com/skpr/operator/pkg/apis/search/v1beta1"
)

func TestSolr(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origSolr := &searchv1beta1.Solr{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: searchv1beta1.SolrSpec{
			Core:    "core1",
			Version: "7",
			Resources: searchv1beta1.SolrSpecResources{
				CPU: searchv1beta1.SolrSpecResourcesCPU{
					Request: resource.MustParse("50m"),
					Limit:   resource.MustParse("500m"),
				},
				Memory:  resource.MustParse("256Mi"),
				Storage: resource.MustParse("10Gi"),
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origSolr, Solr(parent, origSolr.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "Solr result is created")

	newSolr := &searchv1beta1.Solr{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: searchv1beta1.SolrSpec{
			Core:    "core2",
			Version: "7",
			Resources: searchv1beta1.SolrSpecResources{
				CPU: searchv1beta1.SolrSpecResourcesCPU{
					Request: resource.MustParse("50m"),
					Limit:   resource.MustParse("500m"),
				},
				Memory:  resource.MustParse("256Mi"),
				Storage: resource.MustParse("10Gi"),
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origSolr, Solr(parent, newSolr.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Solr result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origSolr, Solr(parent, newSolr.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Solr result is unchanged")
}
