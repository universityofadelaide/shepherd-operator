package backup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/apis"
	extensionv1 "gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/apis/extension/v1"
)

func TestReconcile(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	instance := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: extensionv1.BackupSpec{},
	}

	// Query which will be used to find our Drupal object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &ReconcileBackup{
		Client: fake.NewFakeClient(instance),
		scheme: scheme.Scheme,
	}

	_, err := rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	found := &extensionv1.Backup{}
	err = rd.Client.Get(context.TODO(), query, found)
	assert.Nil(t, err)

}
