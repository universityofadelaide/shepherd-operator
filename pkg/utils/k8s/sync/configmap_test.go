// +build unit

package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestConfigMap(t *testing.T) {
	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Data: map[string]string{
			"1": "2",
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origConfigMap, ConfigMap(parent, origConfigMap.Data, origConfigMap.BinaryData, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "ConfigMap result is created")

	newConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Data: map[string]string{
			"1": "2",
			"3": "4",
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origConfigMap, ConfigMap(parent, newConfigMap.Data, newConfigMap.BinaryData, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "ConfigMap result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origConfigMap, ConfigMap(parent, newConfigMap.Data, newConfigMap.BinaryData, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "ConfigMap result is unchanged")
}
