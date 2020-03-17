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

func TestSecret(t *testing.T) {
	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Data: map[string][]byte{
			"1": []byte("2"),
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origSecret, Secret(parent, origSecret.Data, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "Secret result is created")

	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Data: map[string][]byte{
			"1": []byte("2"),
			"3": []byte("4"),
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origSecret, Secret(parent, newSecret.Data, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Secret result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origSecret, Secret(parent, newSecret.Data, true, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Secret result is unchanged")
}
