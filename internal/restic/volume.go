package restic

import (
	corev1 "k8s.io/api/core/v1"
)

// AttachVolume will add the Restic secrets volume to a Pod.
func AttachVolume(volumes []corev1.Volume) []corev1.Volume {
	mode := corev1.ConfigMapVolumeSourceDefaultMode

	volumes = append(volumes, corev1.Volume{
		Name: VolumeSecrets,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				DefaultMode: &mode,
				SecretName:  ResticSecretPasswordName,
			},
		},
	})

	return volumes
}
