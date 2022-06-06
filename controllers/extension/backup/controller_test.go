package backup

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	extensionv1.AddToScheme(scheme.Scheme)
	batchv1.AddToScheme(scheme.Scheme)

	instance := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: extensionv1.BackupSpec{},
	}

	// Query which will be used to find our Backup object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &Reconciler{
		Client:   fake.NewFakeClient(instance),
		Scheme:   scheme.Scheme,
		Recorder: mockevents.New(),
	}

	_, err := rd.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	found := &extensionv1.Backup{}
	err = rd.Client.Get(context.TODO(), query, found)
	assert.Nil(t, err)
}

func TestReconcileDelete(t *testing.T) {
	extensionv1.AddToScheme(scheme.Scheme)

	deletionTimestamp := &metav1.Time{}
	deletionTimestamp.Time = time.Now()
	instance := &extensionv1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test",
			Namespace:         corev1.NamespaceDefault,
			DeletionTimestamp: deletionTimestamp,
			Finalizers:        []string{Finalizer},
		},
		Spec: extensionv1.BackupSpec{},
		Status: extensionv1.BackupStatus{
			ResticID: "test",
		},
	}

	// Query which will be used to find our Backup object.
	backupQuery := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}
	rd := &Reconciler{
		Client:   fake.NewFakeClient(instance),
		Scheme:   scheme.Scheme,
		Recorder: mockevents.New(),
	}

	_, err := rd.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: backupQuery,
	})
	assert.Nil(t, err)

	// Query which will be used to find our finalizer job object.
	jobName := fmt.Sprintf("restic-delete-%s", backupQuery.Name)
	jobQuery := types.NamespacedName{
		Name:      jobName,
		Namespace: backupQuery.Namespace,
	}
	found := &batchv1.Job{}
	err = rd.Client.Get(context.TODO(), jobQuery, found)
	assert.Nil(t, err)
	assert.Equal(t, jobName, found.Name, "restic delete job found")
}
