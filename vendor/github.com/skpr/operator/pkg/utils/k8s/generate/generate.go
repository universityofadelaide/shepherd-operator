package generate

import (
	corev1 "k8s.io/api/core/v1"
)

// VolumeConfigMap returns a Volume which is used to mount a ConfigMap.
func VolumeConfigMap(name, configmap string) corev1.Volume {
	mode := corev1.ConfigMapVolumeSourceDefaultMode

	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				DefaultMode: &mode,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configmap,
				},
			},
		},
	}
}

// VolumeSecret returns a Volume which is used to mount a Secret.
func VolumeSecret(name, secret string) corev1.Volume {
	mode := corev1.ConfigMapVolumeSourceDefaultMode

	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				DefaultMode: &mode,
				SecretName:  secret,
			},
		},
	}
}

// VolumeClaim returns a Volume which is used to mount a PersistentVolumeClaim.
func VolumeClaim(name, pvc string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc,
			},
		},
	}
}

// Mount a Volume.
func Mount(name, path string, readonly bool) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		MountPath: path,
		ReadOnly:  readonly,
	}
}

// EnvVar which consistens of a key/value pair.
func EnvVar(name, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  name,
		Value: value,
	}
}

// EnvVarConfigMap exposes a ConfigMap key/value as an environment variable.
func EnvVarConfigMap(name, key string, ref *corev1.ConfigMap, optional bool) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key:      key,
				Optional: &optional,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: ref.ObjectMeta.Name,
				},
			},
		},
	}
}

// EnvVarSecret exposes a Secret key/value as an environment variable.
func EnvVarSecret(name, key string, ref *corev1.Secret, optional bool) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key:      key,
				Optional: &optional,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: ref.ObjectMeta.Name,
				},
			},
		},
	}
}
