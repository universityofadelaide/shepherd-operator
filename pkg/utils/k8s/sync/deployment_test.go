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

func TestDeployment(t *testing.T) {
	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	var replicas int32 = 1

	origDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"foo": "bar",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           "foo/bar:0.0.1",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					// The below are fields which need to be set so we can perform an "deep equal"
					// without always having difference.
					SecurityContext: &corev1.PodSecurityContext{},
					SchedulerName:   corev1.DefaultSchedulerName,
					DNSPolicy:       corev1.DNSClusterFirst,
					RestartPolicy:   corev1.RestartPolicyAlways,
				},
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origDeployment, Deployment(parent, origDeployment.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "Deployment result is created")

	newDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"foo": "bar",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           "foo/bar:0.0.2",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					// The below are fields which need to be set so we can perform an "deep equal"
					// without always having difference.
					SecurityContext: &corev1.PodSecurityContext{},
					SchedulerName:   corev1.DefaultSchedulerName,
					DNSPolicy:       corev1.DNSClusterFirst,
					RestartPolicy:   corev1.RestartPolicyAlways,
				},
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDeployment, Deployment(parent, newDeployment.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Deployment result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDeployment, Deployment(parent, newDeployment.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Deployment result is unchanged")
}
