package restic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

// VolumeSecrets identifier used for Restic secret.
const VolumeSecrets = "restic-secrets"

// AttachVolume will add the Restic secrets volume to a Pod.
func AttachVolume(volumes []corev1.Volume, backup *extensionsv1beta1.Backup) []corev1.Volume {
	mode := corev1.ConfigMapVolumeSourceDefaultMode

	volumes = append(volumes, corev1.Volume{
		Name: VolumeSecrets,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				DefaultMode: &mode,
				SecretName:  fmt.Sprintf("%s-%s", Prefix, backup.ObjectMeta.Name),
			},
		},
	})

	return volumes
}
