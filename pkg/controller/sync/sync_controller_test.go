package sync

import (
	"context"
	"testing"

	osv1 "github.com/openshift/api/apps/v1"
	osclient "github.com/openshift/client-go/apps/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/universityofadelaide/shepherd-operator/pkg/apis"
	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

func TestReconcile(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	instance := &extensionv1.Sync{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: extensionv1.SyncSpec{
			Site:       "2",
			BackupEnv:  "3",
			RestoreEnv: "4",
			BackupSpec: extensionv1.BackupSpec{
				Volumes: map[string]shpmetav1.SpecVolume{
					"foo": {
						ClaimName: "bar",
					},
				},
				MySQL: map[string]shpmetav1.SpecMySQL{
					"foo": {
						Secret: shpmetav1.SpecMySQLSecret{
							Name: "bar",
							Keys: shpmetav1.SpecMySQLSecretKeys{
								Username: "test1",
								Password: "test2",
								Database: "test3",
								Hostname: "test4",
								Port:     "test5",
							},
						},
					},
				},
			},
			RestoreSpec: extensionv1.BackupSpec{
				Volumes: map[string]shpmetav1.SpecVolume{
					"foo2": {
						ClaimName: "bar2",
					},
				},
				MySQL: map[string]shpmetav1.SpecMySQL{
					"foo2": {
						Secret: shpmetav1.SpecMySQLSecret{
							Name: "bar2",
							Keys: shpmetav1.SpecMySQLSecretKeys{
								Username: "test11",
								Password: "test22",
								Database: "test33",
								Hostname: "test44",
								Port:     "test55",
							},
						},
					},
				},
			},
		},
	}
	deploymentconfig := &osv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-4",
			Namespace: corev1.NamespaceDefault,
		},
		Status: osv1.DeploymentConfigStatus{
			Conditions: []osv1.DeploymentCondition{
				{
					Type:   osv1.DeploymentAvailable,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}

	// Query which will be used to find our Sync object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &ReconcileSync{
		Client:   fake.NewFakeClient(instance),
		OsClient: osclient.NewSimpleClientset(deploymentconfig).AppsV1(),
		scheme:   scheme.Scheme,
	}

	_, err := rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	sync := &extensionv1.Sync{}
	err = rd.Client.Get(context.TODO(), query, sync)
	assert.Nil(t, err)

	backupName := "sync-test-backup"
	backupQuery := types.NamespacedName{
		Name:      backupName,
		Namespace: instance.ObjectMeta.Namespace,
	}
	backup := &extensionv1.Backup{}
	err = rd.Client.Get(context.TODO(), backupQuery, backup)
	assert.Nil(t, err)
	assert.Equal(t, backupName, backup.Name)
	assert.Equal(t, backup.Spec.Volumes["foo"].ClaimName, "bar")

	restoreName := "sync-test-restore"
	restoreQuery := types.NamespacedName{
		Name:      restoreName,
		Namespace: instance.ObjectMeta.Namespace,
	}
	restore := &extensionv1.Restore{}
	err = rd.Client.Get(context.TODO(), restoreQuery, restore)
	assert.Nil(t, err)
	assert.Equal(t, restoreName, restore.Name)
	assert.Equal(t, restore.Spec.Volumes["foo2"].ClaimName, "bar2")
}
