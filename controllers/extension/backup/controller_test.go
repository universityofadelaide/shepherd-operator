package backup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	mockevents "github.com/universityofadelaide/shepherd-operator/internal/events/mock"
)

func TestReconcile(t *testing.T) {
	err := extensionv1.AddToScheme(scheme.Scheme)
	assert.Nil(t, err)

	err = corev1.AddToScheme(scheme.Scheme)
	assert.Nil(t, err)

	instance := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				"shp_namespace": "test",
			},
		},
		Spec: extensionv1.BackupSpec{},
	}

	// Query which will be used to find our Backup object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &Reconciler{
		Client:    fake.NewClientBuilder().WithObjects(instance).Build(),
		Scheme:    scheme.Scheme,
		Recorder:  mockevents.New(),
		ClientSet: clientfake.NewSimpleClientset(),
		Params: Params{
			ResourceRequirements: corev1.ResourceRequirements{},
			WorkingDir:           "/tmp",
			MySQL: MySQL{
				Image: "mysql:latest",
			},
			AWS: AWS{
				BucketName:     "test",
				Image:          "aws-cli:latest",
				FieldKeyID:     "aws.key.id",
				FieldAccessKey: "aws.access.key",
				Region:         "ap-southeast-2",
			},
			FilterByLabelAndValue: FilterByLabelAndValue{
				Key:   "shp_namespace",
				Value: "test",
			},
		},
	}

	_, err = rd.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	list, err := rd.ClientSet.CoreV1().Pods(query.Namespace).List(context.TODO(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list.Items))
}
