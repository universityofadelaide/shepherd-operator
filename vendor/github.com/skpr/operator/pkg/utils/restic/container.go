package restic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	// EnvAWSAccessKeyID for anthentication.
	EnvAWSAccessKeyID = "AWS_ACCESS_KEY_ID"
	// EnvAWSSecretAccessKey for anthentication.
	EnvAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	// EnvResticRepository for Restic configuration.
	EnvResticRepository = "RESTIC_REPOSITORY"
	// EnvResticPasswordFile for Restic configuration.
	EnvResticPasswordFile = "RESTIC_PASSWORD_FILE"

	// SecretDir defines the directory where secrets are mounted.
	SecretDir = "/etc/restic"
)

// WrapContainer with the information required to interact with Restic.
func WrapContainer(container corev1.Container, key, secret, bucket, repository, passwordFile string) corev1.Container {
	envs := []corev1.EnvVar{
		{
			Name:  EnvAWSAccessKeyID,
			Value: key,
		},
		{
			Name:  EnvAWSSecretAccessKey,
			Value: secret,
		},
		{
			Name:  EnvResticRepository,
			Value: fmt.Sprintf("%s/%s", bucket, repository),
		},
		{
			Name:  EnvResticPasswordFile,
			Value: fmt.Sprintf("%s/%s", SecretDir, passwordFile),
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
