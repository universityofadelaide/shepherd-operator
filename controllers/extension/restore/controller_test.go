package restore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	mockevents "github.com/universityofadelaide/shepherd-operator/internal/events/mock"
)

func TestReconcile(t *testing.T) {
	err := extensionv1.AddToScheme(scheme.Scheme)
	assert.Nil(t, err)

	err = batchv1.AddToScheme(scheme.Scheme)
	assert.Nil(t, err)

	instance := &extensionv1.Restore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: extensionv1.RestoreSpec{},
	}

	// Query which will be used to find our Restore object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &Reconciler{
		Client:   fake.NewClientBuilder().WithObjects(instance).Build(),
		Scheme:   scheme.Scheme,
		Recorder: mockevents.New(),
		Params: Params{
			ResourceRequirements: corev1.ResourceRequirements{},
			WorkingDir:           "/tmp",
			MySQL: MySQL{
				Image: "mysql:latest",
			},
			AWS: AWS{
				BucketName:     "test",
				Image:          "aws-cli:latest",
				SecretName:     "backup",
				FieldKeyID:     "aws.key.id",
				FieldAccessKey: "aws.access.key",
				Region:         "ap-southeast-2",
			},
		},
	}

	_, err = rd.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	found := &extensionv1.Restore{}
	err = rd.Client.Get(context.TODO(), query, found)
	assert.Nil(t, err)
}
