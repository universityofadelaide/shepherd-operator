package image

import (
	"fmt"
	"path/filepath"

	mtkenv "github.com/skpr/mtk/dump/pkg/envar"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	awscredentials "github.com/skpr/operator/pkg/utils/aws/credentials"
)

const (
	// Prefix which is applied to objects created.
	Prefix = "mysql-image"

	// RulesFile used when dumping the database.
	RulesFile = "mtk.yml"

	// ContainerMTKDump for querying the "dump" container.
	ContainerMTKDump = "dump"
	// ContainerKaniko for querying the "build" container.
	ContainerKaniko = "build"

	// MountWorkspaceName for mounting the workspace.
	MountWorkspaceName = "workspace"
	// MountWorkspacePath for mounting the workspace.
	MountWorkspacePath = "/workspace"

	// MountConfigName for mounting the config.
	MountConfigName = "config"
	// MountConfigPath for mounting the config.
	MountConfigPath = "/config"

	// MountDockerName for mounting Docker config.
	MountDockerName = "docker"
	// MountDockerPath for mounting Docker config.
	MountDockerPath = "/kaniko/.docker"
)

// GenerateParams are passed into the Generate function.
type GenerateParams struct {
	Dump   Container
	Build  Container
	Docker GenerateParamsDocker
	AWS    GenerateParamsAWS
}

// GenerateParamsDocker used to Generate a Pod.
type GenerateParamsDocker struct {
	// ConfigMap which contains the value.
	ConfigMap string
}

// GenerateParamsAWS used to Generate a Pod.
type GenerateParamsAWS struct {
	// Secret which contains
	Secret string
	// ConfigMap key which contains AWS_ACCESS_KEY_ID.
	KeyID string
	// KeyAccess key which contains AWS_SECRET_ACCESS_KEY.
	AccessKey string
}

// Generate required objects for building a MySQL image.
func Generate(image *extensionsv1beta1.Image, params GenerateParams) (*corev1.Pod, *corev1.ConfigMap, error) {
	dumpResources, err := params.Dump.ResourceRequirements()
	if err != nil {
		return nil, nil, err
	}

	buildResources, err := params.Build.ResourceRequirements()
	if err != nil {
		return nil, nil, err
	}

	rules, err := yaml.Marshal(&image.Spec.Rules)
	if err != nil {
		return nil, nil, err
	}

	metadata := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", Prefix, image.ObjectMeta.Name),
		Namespace: image.ObjectMeta.Namespace,
	}

	configmap := &corev1.ConfigMap{
		ObjectMeta: metadata,
		Data: map[string]string{
			RulesFile: string(rules),
		},
	}

	volumes := []corev1.Volume{
		{
			Name: MountConfigName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configmap.Name,
					},
				},
			},
		},
		{
			Name: MountWorkspaceName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	mounts := []corev1.VolumeMount{
		{
			Name:      MountConfigName,
			MountPath: MountConfigPath,
		},
		{
			Name:      MountWorkspaceName,
			MountPath: MountWorkspacePath,
		},
	}

	if params.Docker.ConfigMap != "" {
		volumes = append(volumes, corev1.Volume{
			Name: MountDockerName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: params.Docker.ConfigMap,
					},
				},
			},
		})

		mounts = append(mounts, corev1.VolumeMount{
			Name:      MountDockerName,
			MountPath: MountDockerPath,
		})
	}

	var kanikoEnv []corev1.EnvVar

	if params.AWS.Secret != "" {
		if params.AWS.KeyID != "" {
			kanikoEnv = append(kanikoEnv, corev1.EnvVar{
				Name: awscredentials.AccessKeyID,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: params.AWS.Secret,
						},
						Key: params.AWS.KeyID,
					},
				},
			})
		}

		if params.AWS.AccessKey != "" {
			kanikoEnv = append(kanikoEnv, corev1.EnvVar{
				Name: awscredentials.SecretAccessKey,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: params.AWS.Secret,
						},
						Key: params.AWS.AccessKey,
					},
				},
			})
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metadata,
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			InitContainers: []corev1.Container{
				{
					Name:      ContainerMTKDump,
					Image:     params.Dump.Image,
					Resources: dumpResources,
					Env: []corev1.EnvVar{
						{
							Name:  mtkenv.Config,
							Value: filepath.Join(MountWorkspacePath, RulesFile),
						},
						{
							Name: mtkenv.Database,
							ValueFrom: &corev1.EnvVarSource{
								ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: image.Spec.Connection.ConfigMap.Name,
									},
									Key: image.Spec.Connection.ConfigMap.Keys.Database,
								},
							},
						},
						{
							Name: mtkenv.Hostname,
							ValueFrom: &corev1.EnvVarSource{
								ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: image.Spec.Connection.ConfigMap.Name,
									},
									Key: image.Spec.Connection.ConfigMap.Keys.Hostname,
								},
							},
						},
						{
							Name: mtkenv.Port,
							ValueFrom: &corev1.EnvVarSource{
								ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: image.Spec.Connection.ConfigMap.Name,
									},
									Key: image.Spec.Connection.ConfigMap.Keys.Port,
								},
							},
						},
						{
							Name: mtkenv.Username,
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: image.Spec.Connection.Secret.Name,
									},
									Key: image.Spec.Connection.Secret.Keys.Username,
								},
							},
						},
						{
							Name: mtkenv.Password,
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: image.Spec.Connection.Secret.Name,
									},
									Key: image.Spec.Connection.Secret.Keys.Password,
								},
							},
						},
					},
					Command: []string{
						"/bin/bash", "-c",
					},
					Args: []string{
						fmt.Sprintf("mtk-dump > %s/db.sql", MountWorkspacePath),
					},
					VolumeMounts: mounts,
				},
			},
			Containers: []corev1.Container{
				{
					Name:      ContainerKaniko,
					Image:     params.Build.Image,
					Resources: buildResources,
					Env:       kanikoEnv,
					Args: append([]string{
						"--context", MountWorkspacePath,
						"--dockerfile", "/Dockerfile", // @todo, This should be managed by the image.
						"--single-snapshot", "", // @todo, This should be managed by the image.
						"--verbosity", "fatal", // @todo, This should be managed by the image.
					}, formatDestination(image.Spec.Destinations)...),
					VolumeMounts: mounts,
				},
			},
			Volumes: volumes,
		},
	}

	return pod, configmap, nil
}
