package restic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
)

const (
	// EnvMySQLHostname for MySQL connection.
	EnvMySQLHostname = "MYSQL_HOSTNAME"
	// EnvMySQLDatabase for MySQL connection.
	EnvMySQLDatabase = "MYSQL_DATABASE"
	// EnvMySQLPort for MySQL connection.
	EnvMySQLPort = "MYSQL_PORT"
	// EnvMySQLUsername for MySQL connection.
	EnvMySQLUsername = "MYSQL_USERNAME"
	// EnvMySQLPassword for MySQL connection.
	EnvMySQLPassword = "MYSQL_PASSWORD"

	// VolumeMySQL identifier for mysql storage.
	VolumeMySQL = "restic-mysql"
	// VolumePublic identifier for public storage.
	VolumePublic = "restic-public"
	// VolumePrivate identifier for private storage.
	VolumePrivate = "restic-private"
)

// PodSpecParams which are passed into the PodSpec function.
type PodSpecParams struct {
	Bucket      string
	KeyID       string
	AccessKey   string
	CPU         string
	Memory      string
	ResticImage string
	MySQLImage  string
	WorkingDir  string
	Tags        []string
}

// PodSpec defines how a backup can be executed using a Pod.
func PodSpec(backup *extensionsv1beta1.Backup, params PodSpecParams) (corev1.PodSpec, error) {
	cpu, err := resource.ParseQuantity(params.CPU)
	if err != nil {
		return corev1.PodSpec{}, err
	}

	memory, err := resource.ParseQuantity(params.Memory)
	if err != nil {
		return corev1.PodSpec{}, err
	}

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
	}

	resticInit := WrapContainer(corev1.Container{
		Name:      "restic-init",
		Image:     params.ResticImage,
		Resources: resources,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			// Init will return an exit code of 1 if the repository alread exists.
			// If this failed for a non "already exists" error then we will see it
			// in the main containers "restic backup" execution.
			"restic init || true",
		},
	}, params.KeyID, params.AccessKey, params.Bucket, backup)

	resticBackup := corev1.Container{
		Name:       "restic-backup",
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: params.WorkingDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			fmt.Sprintf("restic --verbose %s backup .", formatTags(params.Tags)),
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      VolumeMySQL,
				MountPath: fmt.Sprintf("%s/mysql", params.WorkingDir),
			},
		},
	}

	for volumeName := range backup.Spec.Volumes {
		resticBackup.VolumeMounts = append(resticBackup.VolumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("volume-%s", volumeName),
			MountPath: fmt.Sprintf("%s/volume/%s", params.WorkingDir, volumeName),
			ReadOnly:  true,
		})
	}

	resticBackup = WrapContainer(resticBackup, params.KeyID, params.AccessKey, params.Bucket, backup)

	spec := corev1.PodSpec{
		RestartPolicy: corev1.RestartPolicyNever,
		InitContainers: []corev1.Container{
			resticInit,
		},
		Containers: []corev1.Container{
			resticBackup,
		},
		Volumes: AttachVolume([]corev1.Volume{
			{
				Name: VolumeMySQL,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: corev1.StorageMediumDefault,
					},
				},
			},
		}, backup),
	}

	for volumeName, volumeSpec := range backup.Spec.Volumes {
		spec.Volumes = append(spec.Volumes, corev1.Volume{
			Name: fmt.Sprintf("volume-%s", volumeName),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeSpec.ClaimName,
				},
			},
		})
	}

	for mysqlName, mysqlStatus := range backup.Spec.MySQL {
		spec.InitContainers = append(spec.InitContainers, corev1.Container{
			Name:      fmt.Sprintf("mysql-%s", mysqlName),
			Image:     params.MySQLImage,
			Resources: resources,
			Env: []corev1.EnvVar{
				{
					Name: EnvMySQLHostname,
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: mysqlStatus.ConfigMap.Name,
							},
							Key: mysqlStatus.ConfigMap.Keys.Hostname,
						},
					},
				},
				{
					Name: EnvMySQLDatabase,
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: mysqlStatus.ConfigMap.Name,
							},
							Key: mysqlStatus.ConfigMap.Keys.Database,
						},
					},
				},
				{
					Name: EnvMySQLPort,
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: mysqlStatus.ConfigMap.Name,
							},
							Key: mysqlStatus.ConfigMap.Keys.Port,
						},
					},
				},
				{
					Name: EnvMySQLUsername,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: mysqlStatus.Secret.Name,
							},
							Key: mysqlStatus.Secret.Keys.Username,
						},
					},
				},
				{
					Name: EnvMySQLPassword,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: mysqlStatus.Secret.Name,
							},
							Key: mysqlStatus.Secret.Keys.Password,
						},
					},
				},
			},
			WorkingDir: params.WorkingDir,
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				fmt.Sprintf("mysqldump --single-transaction --host=\"$MYSQL_HOSTNAME\" --user=\"$MYSQL_USERNAME\" --password=\"$MYSQL_PASSWORD\" \"$MYSQL_DATABASE\" > \"mysql/%s.sql\"", mysqlName),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", params.WorkingDir),
				},
			},
		})
	}

	return spec, nil
}
