package restic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
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

	// ResticRepoDir defines the directory to mount the restic repository to.
	ResticRepoDir = "/srv/backups"
)

// WrapContainer with the information required to interact with Restic.
func WrapContainer(container corev1.Container, siteId, namespace string) corev1.Container {
	envs := []corev1.EnvVar{
		{
			Name:  EnvResticRepository,
			Value: fmt.Sprintf("%s/%s/%s", ResticRepoDir, namespace, siteId),
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

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      VolumeRepository,
		MountPath: ResticRepoDir,
	})

	return container
}
