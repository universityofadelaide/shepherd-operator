// +build unit

package sync

import (
	"context"
	"testing"

	"github.com/skpr/operator/pkg/apis"
	"github.com/skpr/operator/pkg/mysql"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

func TestDatabase(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origDatabase := &mysqlv1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: mysqlv1beta1.DatabaseSpec{
			Provisioner: "test",
			Privileges: []string{
				mysql.PrivilegeSelect,
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origDatabase, Database(parent, origDatabase.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "Database result is created")

	newDatabase := &mysqlv1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: mysqlv1beta1.DatabaseSpec{
			Provisioner: "test",
			Privileges: []string{
				mysql.PrivilegeDrop,
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDatabase, Database(parent, newDatabase.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Database result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDatabase, Database(parent, newDatabase.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Database result is unchanged")
}
