package restic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	extensionv1 "gitlab.adelaide.edu.au/web-team/shepherd-operator/pkg/apis/extension/v1"
)

const (
	// EnvResticRepository for Restic configuration.
	EnvResticRepository = "RESTIC_REPOSITORY"
	// EnvResticPasswordFile for Restic configuration.
	EnvResticPasswordFile = "RESTIC_PASSWORD_FILE"

	// ResticPassword identifier for loading the restic password.
	ResticPassword = "password"

	// SecretDir defines the directory where secrets are mounted.
	SecretDir = "/etc/restic"
)

// WrapContainer with the information required to interact with Restic.
func WrapContainer(container corev1.Container, siteId string, backup *extensionv1.Backup) corev1.Container {
	envs := []corev1.EnvVar{
		{
			Name:  EnvResticRepository,
			Value: fmt.Sprintf("%s/%s/%s", siteId, backup.ObjectMeta.Namespace, backup.ObjectMeta.Name),
		},
		{
			Name:  EnvResticPasswordFile,
			Value: fmt.Sprintf("%s/%s", SecretDir, ResticPassword),
		},
	}

	container.Env = append(container.Env, envs...)

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      VolumeSecrets,
		MountPath: SecretDir,
		ReadOnly:  true,
	})

	return container
}
