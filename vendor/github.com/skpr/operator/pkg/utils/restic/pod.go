package restic

import (
	"fmt"
	"strings"

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
)

// PodSpecParams which are passed into the PodSpecBackup function.
type PodSpecParams struct {
	Bucket      string
	KeyID       string
	AccessKey   string
	CPU         string
	Memory      string
	ResticImage string
	MySQLImage  string
	WorkingDir  string
}

// PodSpecBackup defines how a backup can be executed using a Pod.
func PodSpecBackup(backup *extensionsv1beta1.Backup, params PodSpecParams) (corev1.PodSpec, error) {
	resources, err := resourcesFromParams(params)
	if err != nil {
		return corev1.PodSpec{}, err
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
	}, params.KeyID, params.AccessKey, params.Bucket, backup.ObjectMeta.Namespace, backup.Spec.Secret.Key)

	resticBackup := corev1.Container{
		Name:       ResticBackupContainerName,
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: params.WorkingDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			fmt.Sprintf("restic --verbose %s backup .", formatTags(backup.Spec.Tags)),
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

	resticBackup = WrapContainer(resticBackup, params.KeyID, params.AccessKey, params.Bucket, backup.ObjectMeta.Namespace, backup.Spec.Secret.Key)

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
		}, backup.Spec.Secret.Name),
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
			Name:       fmt.Sprintf("mysql-%s", mysqlName),
			Image:      params.MySQLImage,
			Resources:  resources,
			Env:        mysqlEnvVars(mysqlStatus),
			WorkingDir: params.WorkingDir,
			Command: []string{
				"bash",
				"-c",
			},
			Args: []string{
				fmt.Sprintf("database-backup > mysql/%s.sql", mysqlName),
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

// PodSpecRestore defines how a restore can be executed using a Pod.
func PodSpecRestore(restore *extensionsv1beta1.Restore, resticID string, params PodSpecParams) (corev1.PodSpec, error) {
	resources, err := resourcesFromParams(params)
	if err != nil {
		return corev1.PodSpec{}, err
	}

	var (
		initContainers []corev1.Container
		containers     []corev1.Container
	)

	for mysqlName, mysqlStatus := range restore.Spec.MySQL {
		// InitContainers which restores each sql file to an emptydir volume.
		initContainers = append(initContainers, corev1.Container{
			Name:       fmt.Sprintf("restic-restore-%s", mysqlName),
			Image:      params.ResticImage,
			Resources:  resources,
			WorkingDir: params.WorkingDir,
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				fmt.Sprintf("restic dump --quiet %s /%[2]s > ./%[2]s", resticID, fmt.Sprintf("mysql/%s.sql", mysqlName)),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", params.WorkingDir),
				},
			},
		})

		// Containers which will restore the database.
		containers = append(containers, corev1.Container{
			Name:       fmt.Sprintf("restic-import-%s", mysqlName),
			Image:      params.MySQLImage,
			Resources:  resources,
			WorkingDir: params.WorkingDir,
			Command: []string{
				"database-restore",
			},
			Args: []string{
				fmt.Sprintf("mysql/%s.sql", mysqlName),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", params.WorkingDir),
				},
			},
			Env: mysqlEnvVars(mysqlStatus),
		})
	}

	// Volume definitions for the pod.
	specVolumes := AttachVolume([]corev1.Volume{
		{
			Name: VolumeMySQL,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}, restore.Spec.Secret.Name)

	// Gather volumes from the restore that need to be mounted.
	var (
		resticRestoreVolumeMounts []corev1.VolumeMount
		resticVolumeIncludeArgs   []string
	)
	for volumeName, volumeSpec := range restore.Spec.Volumes {
		resticVolumeIncludeArgs = append(resticVolumeIncludeArgs, fmt.Sprintf("--include /volume/%s", volumeName))
		volumeName := fmt.Sprintf("volume-%s", volumeName)
		resticRestoreVolumeMounts = append(resticRestoreVolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: fmt.Sprintf("%s/volume/%s", params.WorkingDir, volumeName),
			ReadOnly:  false,
		})

		specVolumes = append(specVolumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeSpec.ClaimName,
				},
			},
		})
	}

	// Container which restores volumes.
	containers = append(containers, corev1.Container{
		Name:       "restic-restore-volumes",
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: params.WorkingDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			fmt.Sprintf("restic restore %s --target . %s", resticID, strings.Join(resticVolumeIncludeArgs, " ")),
		},
		VolumeMounts: resticRestoreVolumeMounts,
	})

	for i := range initContainers {
		initContainers[i] = WrapContainer(initContainers[i], params.KeyID, params.AccessKey, params.Bucket, restore.ObjectMeta.Namespace, restore.Spec.Secret.Key)
	}
	for i := range containers {
		containers[i] = WrapContainer(containers[i], params.KeyID, params.AccessKey, params.Bucket, restore.ObjectMeta.Namespace, restore.Spec.Secret.Key)
	}
	return corev1.PodSpec{
		RestartPolicy:  corev1.RestartPolicyNever,
		InitContainers: initContainers,
		Containers:     containers,
		Volumes:        specVolumes,
	}, nil
}

func resourcesFromParams(params PodSpecParams) (corev1.ResourceRequirements, error) {
	cpu, err := resource.ParseQuantity(params.CPU)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}

	memory, err := resource.ParseQuantity(params.Memory)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpu,
			corev1.ResourceMemory: memory,
		},
	}, nil
}

func mysqlEnvVars(mysqlStatus extensionsv1beta1.BackupSpecMySQL) []corev1.EnvVar {
	return []corev1.EnvVar{
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
	}
}
