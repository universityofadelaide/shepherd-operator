/*
Copyright 2019 University of Adelaide.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backupscheduled

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

	"github.com/universityofadelaide/shepherd-operator/pkg/apis"
	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
	clock "github.com/universityofadelaide/shepherd-operator/pkg/utils/clock/mock"
	events "github.com/universityofadelaide/shepherd-operator/pkg/utils/k8s/events/mock"
)

func TestReconcile(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	retentionMaxNumber := 2
	startDeadline := int64(60)
	instance := &extensionv1.BackupScheduled{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				"site": "foo",
			},
		},
		Spec: extensionv1.BackupScheduledSpec{
			Retention: shpmetav1.RetentionSpec{
				MaxNumber: &retentionMaxNumber,
			},
			Schedule: shpmetav1.ScheduledSpec{
				CronTab:                 "0 0 * * *",
				StartingDeadlineSeconds: &startDeadline,
			},
		},
	}

	// Query which will be used to find our BackupScheduled object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	c, err := clock.New("2020-04-02T00:00:30Z")
	assert.Nil(t, err)

	rd := &ReconcileBackupScheduled{
		Client:   fake.NewFakeClient(instance),
		scheme:   scheme.Scheme,
		Clock:    c,
		recorder: &events.Mock{},
	}

	_, err = rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Nil(t, err)

	found := &extensionv1.BackupScheduled{}
	err = rd.Client.Get(context.TODO(), query, found)
	assert.Nil(t, err)
}

func TestReconcileNoLabels(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	instance := &extensionv1.BackupScheduled{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	// Query which will be used to find our BackupScheduled object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &ReconcileBackupScheduled{
		Client: fake.NewFakeClient(instance),
		scheme: scheme.Scheme,
	}

	_, err := rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Error(t, err, "BackupScheduled doesn't have a site label.")
}

func TestReconcileNoSchedule(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	instance := &extensionv1.BackupScheduled{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				"site": "foo",
			},
		},
	}

	// Query which will be used to find our BackupScheduled object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	rd := &ReconcileBackupScheduled{
		Client: fake.NewFakeClient(instance),
		scheme: scheme.Scheme,
	}

	_, err := rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Error(t, err, "BackupScheduled doesn't have a schedule.")
}

func TestReconcileInvalidSchedule(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	instance := &extensionv1.BackupScheduled{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
			Labels: map[string]string{
				"site": "foo",
			},
		},
		Spec: extensionv1.BackupScheduledSpec{
			Schedule: shpmetav1.ScheduledSpec{
				CronTab: "a b * * * * *",
			},
		},
	}

	// Query which will be used to find our BackupScheduled object.
	query := types.NamespacedName{
		Name:      instance.ObjectMeta.Name,
		Namespace: instance.ObjectMeta.Namespace,
	}

	c, err := clock.New("2020-04-02T00:00:03Z")
	assert.Nil(t, err)
	rd := &ReconcileBackupScheduled{
		Client: fake.NewFakeClient(instance),
		scheme: scheme.Scheme,
		Clock:  c,
	}

	_, err = rd.Reconcile(reconcile.Request{
		NamespacedName: query,
	})
	assert.Contains(t, err.Error(), "expected exactly 5 fields, found 7")
}
