// +build unit

package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestCronJob(t *testing.T) {
	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "test",
							Image:           "foo/bar:0.0.1",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"/bin/bash",
								"-c",
							},
							Args: []string{
								"echo 1",
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
					RestartPolicy:   corev1.RestartPolicyNever,
				},
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origJob, Job(parent, origJob.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "CronJob result is created")

	newJob := &batchv1.Job{}
	origJob.DeepCopyInto(newJob)
	newJob.Spec.Template.Spec.Containers[0].Name = "Test 2"

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origJob, Job(parent, newJob.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Job result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origJob, Job(parent, newJob.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Job result is unchanged")
}
