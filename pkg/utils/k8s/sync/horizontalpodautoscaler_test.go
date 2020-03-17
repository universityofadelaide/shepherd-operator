// +build unit

package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestHPA(t *testing.T) {
	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	var (
		min int32 = 1
		max int32 = 1
		cpu int32 = 80
		mem int32 = 80
	)

	origHPA := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "test",
			},
			MinReplicas: &min,
			MaxReplicas: max,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: &cpu,
						},
					},
				},
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceMemory,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: &mem,
						},
					},
				},
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origHPA, HPA(parent, origHPA.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "HorizontalPodAutoscaler result is created")

	newHPA := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "test2",
			},
			MinReplicas: &min,
			MaxReplicas: max,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: &cpu,
						},
					},
				},
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceMemory,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: &mem,
						},
					},
				},
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origHPA, HPA(parent, newHPA.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "HorizontalPodAutoscaler result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origHPA, HPA(parent, newHPA.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "HorizontalPodAutoscaler result is unchanged")
}
