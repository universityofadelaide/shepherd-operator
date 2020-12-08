package restic

import (
	"fmt"
	"github.com/pkg/errors"
	v1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"

	osv1 "github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"

	"github.com/universityofadelaide/shepherd-operator/pkg/utils/helper"
	"strings"
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
	CPU         string
	Memory      string
	ResticImage string
	MySQLImage  string
	WorkingDir  string
	Tags        []string
}

// PodSpecBackup defines how a backup can be executed using a Pod.
func PodSpecBackup(backup *extensionv1.Backup, params PodSpecParams, siteId string) (corev1.PodSpec, error) {
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
	}, siteId, backup.ObjectMeta.Namespace)

	resticBackup := corev1.Container{
		Name:       ResticBackupContainerName,
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: params.WorkingDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			// Backup, excluding any cached twig, css, or js.
			fmt.Sprintf("restic --verbose %s backup . --exclude volume/*/*/php --exclude volume/*/*/css --exclude volume/*/*/js", formatTags(params.Tags)),
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

	resticBackup = WrapContainer(resticBackup, siteId, backup.ObjectMeta.Namespace)

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
			{
				Name: VolumeRepository,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: VolumeRepository,
					},
				},
			},
		}),
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
func PodSpecRestore(restore *extensionv1.Restore, dc *osv1.DeploymentConfig, resticId string, params PodSpecParams, siteId string) (corev1.PodSpec, error) {
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

	var initContainers []corev1.Container
	var containers []corev1.Container

	// InitContainer which restores db to emptydir volume.
	for mysqlName, mysqlStatus := range restore.Spec.MySQL {
		initContainers = append(initContainers, corev1.Container{
			Name:       fmt.Sprintf("restic-restore-%s", mysqlName),
			Image:      params.ResticImage,
			Resources:  resources,
			WorkingDir: params.WorkingDir,
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				helper.TprintfMustParse(
					"restic dump --quiet {{.ResticId}} /{{.SQLPath}} > ./{{.SQLPath}}",
					map[string]interface{}{
						"ResticId": resticId,
						"SQLPath":  fmt.Sprintf("mysql/%s.sql", mysqlName),
					},
				),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", params.WorkingDir),
				},
			},
		})

		initContainers = append(initContainers, corev1.Container{
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

	// Mount restore volumes into volume restore container.
	var resticRestoreVolumeMounts []corev1.VolumeMount
	var resticVolumeIncludeArgs []string
	for volumeName := range restore.Spec.Volumes {
		resticVolumeIncludeArgs = append(resticVolumeIncludeArgs, fmt.Sprintf("--include /volume/%s", volumeName))
		resticRestoreVolumeMounts = append(resticRestoreVolumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("volume-%s", volumeName),
			MountPath: fmt.Sprintf("%s/volume/%s", params.WorkingDir, volumeName),
			ReadOnly:  false,
		})
	}

	// Container which restores volumes.
	initContainers = append(initContainers, corev1.Container{
		Name:       "restic-restore-volumes",
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: params.WorkingDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			helper.TprintfMustParse(
				"restic restore {{.ResticId}} --target . {{.IncludeArgs}}",
				map[string]interface{}{
					"ResticId": resticId,
					// @todo might be able to iterate through an array of volumeNames in the template.
					"IncludeArgs": strings.Join(resticVolumeIncludeArgs, " "),
				},
			),
		},
		VolumeMounts: resticRestoreVolumeMounts,
	})

	// Volume definitions for the pod.
	specVolumes := AttachVolume([]corev1.Volume{
		{
			Name: VolumeMySQL,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: VolumeRepository,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: VolumeRepository,
				},
			},
		},
	})
	// Attach restore volumes to pod.
	for volumeName, volumeSpec := range restore.Spec.Volumes {
		specVolumes = append(specVolumes, corev1.Volume{
			Name: fmt.Sprintf("volume-%s", volumeName),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeSpec.ClaimName,
				},
			},
		})
	}

	dcContainer, err := getWebContainerFromDc(dc)
	if err != nil {
		return corev1.PodSpec{}, err
	}
	dcVolumeMounts := dcContainer.VolumeMounts
	// Add volumes from the deploymentconfig that we don't already have in the restore spec.
	for _, dcVolume := range dc.Spec.Template.Spec.Volumes {
		found := false
		for _, specVolume := range specVolumes {
			if dcVolume.PersistentVolumeClaim != nil && specVolume.PersistentVolumeClaim != nil &&
				dcVolume.PersistentVolumeClaim.ClaimName == specVolume.PersistentVolumeClaim.ClaimName {
				found = true
				// We've found a volume we already have, make sure the volume mount name references the existing volume.
				for i, dcVolumeMount := range dcVolumeMounts {
					if dcVolumeMount.Name == dcVolume.Name {
						dcVolumeMounts[i].Name = specVolume.Name
					}
				}
			}
		}

		if !found {
			specVolumes = append(specVolumes, dcVolume)
		}
	}
	// Container which runs deployment steps.
	containers = append(containers, corev1.Container{
		Name:       "restore-deploy",
		Image:      dcContainer.Image,
		Resources:  resources,
		WorkingDir: WebDirectory,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			helper.TprintfMustParse(
				"export REDIS_ENABLED=0 && export MEMCACHE_ENABLED=0 && drush -r {{.WebDir}}/web cr && drush -r {{.WebDir}}/web -y updb && robo config:import-plus && drush -r {{.WebDir}}/web cr",
				map[string]interface{}{
					"WebDir": WebDirectory,
				},
			),
		},
		Env:          dcContainer.Env,
		VolumeMounts: dcVolumeMounts,
	})

	for i, _ := range initContainers {
		initContainers[i] = WrapContainer(initContainers[i], siteId, restore.ObjectMeta.Namespace)
	}
	spec := corev1.PodSpec{
		RestartPolicy:  corev1.RestartPolicyNever,
		InitContainers: initContainers,
		Containers:     containers,
		Volumes:        specVolumes,
	}

	return spec, nil
}

// PodSpecDelete defines how a snapshot can be forgotten using a Pod.
func PodSpecDelete(resticId, namespace, site string, params PodSpecParams) (corev1.PodSpec, error) {
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

	var containers []corev1.Container
	containers = append(containers, WrapContainer(corev1.Container{
		Name:       fmt.Sprintf("restic-delete-%s", resticId),
		Image:      params.ResticImage,
		Resources:  resources,
		WorkingDir: ResticRepoDir,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			fmt.Sprintf("restic forget --prune %s", resticId),
		},
	}, site, namespace))

	spec := corev1.PodSpec{
		RestartPolicy: corev1.RestartPolicyNever,
		Containers:    containers,
		Volumes: AttachVolume([]corev1.Volume{
			{
				Name: VolumeRepository,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: VolumeRepository,
					},
				},
			},
		}),
	}

	return spec, nil
}

// mysqlEnvVars returns a list of environment variables for a container based on the mysql spec.
func mysqlEnvVars(mysqlStatus v1.SpecMySQL) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: EnvMySQLHostname,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mysqlStatus.Secret.Name,
					},
					Key: mysqlStatus.Secret.Keys.Hostname,
				},
			},
		},
		{
			Name: EnvMySQLDatabase,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mysqlStatus.Secret.Name,
					},
					Key: mysqlStatus.Secret.Keys.Database,
				},
			},
		},
		{
			Name: EnvMySQLPort,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mysqlStatus.Secret.Name,
					},
					Key: mysqlStatus.Secret.Keys.Port,
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

// getWebContainerFromDc loops through a deploymentconfig to find the container with the same name. This is considered
// the web container in shepherd.
func getWebContainerFromDc(dc *osv1.DeploymentConfig) (corev1.Container, error) {
	for _, container := range dc.Spec.Template.Spec.Containers {
		if container.Name == dc.ObjectMeta.Name {
			return container, nil
		}
	}
	return corev1.Container{}, errors.Errorf("web container not found for dc %s", dc.ObjectMeta.Name)
}
