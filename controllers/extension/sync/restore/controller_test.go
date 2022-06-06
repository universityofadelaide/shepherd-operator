package restore

import (
	"context"
	osv1 "github.com/openshift/api/apps/v1"
	osclient "github.com/openshift/client-go/apps/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
)

func TestReconcile(t *testing.T) {
	extensionv1.AddToScheme(scheme.Scheme)

	now := metav1.Now()

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
		Status: extensionv1.SyncStatus{
			Backup: extensionv1.SyncStatusBackup{
				Name:      "sync-test-restore",
				Phase:     shpmetav1.PhaseCompleted,
				StartTime: &now,
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

	rd := &Reconciler{
		Client:    fake.NewFakeClient(instance),
		OpenShift: osclient.NewSimpleClientset(deploymentconfig).AppsV1(),
		Scheme:    scheme.Scheme,
	}

	_, err := rd.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	sync := &extensionv1.Sync{}
	err = rd.Client.Get(context.TODO(), query, sync)
	assert.Nil(t, err)

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
