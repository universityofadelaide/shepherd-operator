// +build unit

package sync

import (
	"context"
	"testing"

	"github.com/skpr/operator/pkg/apis"
	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func TestDrupal(t *testing.T) {
	apis.AddToScheme(scheme.Scheme)

	client := fake.NewFakeClient()

	parent := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
	}

	origDrupal := &appv1beta1.Drupal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appv1beta1.DrupalSpec{
			Nginx: appv1beta1.DrupalSpecNginx{
				Image: "drupal:0.0.1-nginx",
				Port:  80,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
				Autoscaling: appv1beta1.DrupalSpecNginxAutoscaling{
					Trigger: appv1beta1.DrupalSpecNginxAutoscalingTrigger{
						CPU: 80,
					},
					Replicas: appv1beta1.DrupalSpecNginxAutoscalingReplicas{
						Min: 2,
						Max: 4,
					},
				},
				HostAlias: appv1beta1.DrupalSpecNginxHostAlias{
					FPM: "php-fpm",
				},
			},
			FPM: appv1beta1.DrupalSpecFPM{
				Image: "drupal:0.0.1-fpm",
				Port:  9000,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("150m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("150m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				},
				Autoscaling: appv1beta1.DrupalSpecFPMAutoscaling{
					Trigger: appv1beta1.DrupalSpecFPMAutoscalingTrigger{
						CPU: 80,
					},
					Replicas: appv1beta1.DrupalSpecFPMAutoscalingReplicas{
						Min: 4,
						Max: 8,
					},
				},
			},
			Volume: appv1beta1.DrupalSpecVolumes{
				Public: appv1beta1.DrupalSpecVolume{
					Path:   "/data/app/sites/default/files",
					Class:  "standard",
					Amount: "10Gi",
				},
				Private: appv1beta1.DrupalSpecVolume{
					Path:   "/mnt/private",
					Class:  "standard",
					Amount: "10Gi",
				},
				Temporary: appv1beta1.DrupalSpecVolume{
					Path:   "/mnt/temporary",
					Class:  "standard",
					Amount: "10Gi",
				},
			},
			MySQL: map[string]appv1beta1.DrupalSpecMySQL{
				"default": appv1beta1.DrupalSpecMySQL{
					Class: "nonprod",
				},
				"migration": appv1beta1.DrupalSpecMySQL{
					Class: "nonprod",
				},
			},
			Cron: map[string]appv1beta1.DrupalSpecCron{
				"drush": appv1beta1.DrupalSpecCron{
					Image:    "drupal:0.0.1-cli",
					Command:  "drush cron",
					Schedule: "*/5 * * * *",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("150m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("150m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
					Retries:     2,
					KeepSuccess: 10,
					KeepFailed:  10,
				},
			},
			ConfigMap: appv1beta1.DrupalSpecConfigMaps{
				Default: appv1beta1.DrupalSpecConfigMap{
					Path: "/etc/skpr/config/default",
				},
				Override: appv1beta1.DrupalSpecConfigMap{
					Path: "/etc/skpr/config/override",
				},
			},
			Secret: appv1beta1.DrupalSpecSecrets{
				Default: appv1beta1.DrupalSpecSecret{
					Path: "/etc/skpr/secret/default",
				},
				Override: appv1beta1.DrupalSpecSecret{
					Path: "/etc/skpr/secret/override",
				},
			},
			NewRelic: appv1beta1.DrupalSpecNewRelic{
				ConfigMap: appv1beta1.DrupalSpecNewRelicConfigMap{
					Enabled: "newrelic.enabled",
					Name:    "newrelic.name",
				},
				Secret: appv1beta1.DrupalSpecNewRelicSecret{
					License: "newrelic.license",
				},
			},
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), client, origDrupal, Drupal(parent, origDrupal.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultCreated), string(result), "Drupal result is created")

	newDrupal := &appv1beta1.Drupal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appv1beta1.DrupalSpec{
			Nginx: appv1beta1.DrupalSpecNginx{
				Image: "drupal:0.0.2-nginx",
				Port:  80,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
				Autoscaling: appv1beta1.DrupalSpecNginxAutoscaling{
					Trigger: appv1beta1.DrupalSpecNginxAutoscalingTrigger{
						CPU: 80,
					},
					Replicas: appv1beta1.DrupalSpecNginxAutoscalingReplicas{
						Min: 2,
						Max: 4,
					},
				},
				HostAlias: appv1beta1.DrupalSpecNginxHostAlias{
					FPM: "php-fpm",
				},
			},
			FPM: appv1beta1.DrupalSpecFPM{
				Image: "drupal:0.0.2-fpm",
				Port:  9000,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("150m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("150m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
				},
				Autoscaling: appv1beta1.DrupalSpecFPMAutoscaling{
					Trigger: appv1beta1.DrupalSpecFPMAutoscalingTrigger{
						CPU: 80,
					},
					Replicas: appv1beta1.DrupalSpecFPMAutoscalingReplicas{
						Min: 4,
						Max: 8,
					},
				},
			},
			Volume: appv1beta1.DrupalSpecVolumes{
				Public: appv1beta1.DrupalSpecVolume{
					Path:   "/data/app/sites/default/files",
					Class:  "standard",
					Amount: "10Gi",
				},
				Private: appv1beta1.DrupalSpecVolume{
					Path:   "/mnt/private",
					Class:  "standard",
					Amount: "10Gi",
				},
				Temporary: appv1beta1.DrupalSpecVolume{
					Path:   "/mnt/temporary",
					Class:  "standard",
					Amount: "10Gi",
				},
			},
			MySQL: map[string]appv1beta1.DrupalSpecMySQL{
				"default": appv1beta1.DrupalSpecMySQL{
					Class: "nonprod",
				},
				"migration": appv1beta1.DrupalSpecMySQL{
					Class: "nonprod",
				},
			},
			Cron: map[string]appv1beta1.DrupalSpecCron{
				"drush": appv1beta1.DrupalSpecCron{
					Image:    "drupal:0.0.2-cli",
					Command:  "drush cron",
					Schedule: "*/5 * * * *",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("150m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("150m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
					Retries:     2,
					KeepSuccess: 10,
					KeepFailed:  10,
				},
			},
			ConfigMap: appv1beta1.DrupalSpecConfigMaps{
				Default: appv1beta1.DrupalSpecConfigMap{
					Path: "/etc/skpr/config/default",
				},
				Override: appv1beta1.DrupalSpecConfigMap{
					Path: "/etc/skpr/config/override",
				},
			},
			Secret: appv1beta1.DrupalSpecSecrets{
				Default: appv1beta1.DrupalSpecSecret{
					Path: "/etc/skpr/secret/default",
				},
				Override: appv1beta1.DrupalSpecSecret{
					Path: "/etc/skpr/secret/override",
				},
			},
			NewRelic: appv1beta1.DrupalSpecNewRelic{
				ConfigMap: appv1beta1.DrupalSpecNewRelicConfigMap{
					Enabled: "newrelic.enabled",
					Name:    "newrelic.name",
				},
				Secret: appv1beta1.DrupalSpecNewRelicSecret{
					License: "newrelic.license",
				},
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDrupal, Drupal(parent, newDrupal.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultUpdated), string(result), "Drupal result is updated")

	result, err = controllerutil.CreateOrUpdate(context.TODO(), client, origDrupal, Drupal(parent, newDrupal.Spec, scheme.Scheme))
	assert.Nil(t, err)
	assert.Equal(t, string(controllerutil.OperationResultNone), string(result), "Drupal result is unchanged")
}
