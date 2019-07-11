package restic

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	osv1 "github.com/openshift/api/apps/v1"
	extensionv1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/extension/v1"
	metav1_shepherd "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

func TestPodSpecBackup(t *testing.T) {
	var params = PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "test/image",
		MySQLImage:  "test/mysqlimage",
		WorkingDir:  "/home/test",
		Tags:        []string{"tag1"},
	}
	backup := extensionv1.Backup{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-backup",
			Namespace: "test-namespace",
		},
		Spec: extensionv1.BackupSpec{
			Volumes: map[string]metav1_shepherd.SpecVolume{
				"volume1": {
					ClaimName: "claim-volume1",
				},
			},
			MySQL: map[string]metav1_shepherd.SpecMySQL{
				"mysql1": {
					Secret: metav1_shepherd.SpecMySQLSecret{
						Name: "secret1",
						Keys: metav1_shepherd.SpecMySQLSecretKeys{
							Username: "mysql-user",
							Password: "mysql-pass",
							Database: "mysql-db",
							Hostname: "mysql-host",
							Port:     "mysql-port",
						},
					},
				},
			},
		},
	}
	cpu, _ := resource.ParseQuantity(params.CPU)
	memory, _ := resource.ParseQuantity(params.Memory)
	mode := corev1.ConfigMapVolumeSourceDefaultMode
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
	expected := corev1.PodSpec{
		RestartPolicy: corev1.RestartPolicyNever,
		InitContainers: []corev1.Container{
			{
				Name:      "restic-init",
				Image:     "test/image",
				Resources: resources,
				Command: []string{
					"/bin/sh", "-c",
				},
				Args: []string{
					"restic init || true",
				},
				Env: []corev1.EnvVar{
					{
						Name:  EnvResticRepository,
						Value: "/srv/backups/test-namespace/test-site-id",
					},
					{
						Name:  EnvResticPasswordFile,
						Value: "/etc/restic/password",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      VolumeSecrets,
						MountPath: SecretDir,
						ReadOnly:  true,
					},
					{
						Name:      VolumeRepository,
						MountPath: ResticRepoDir,
					},
				},
			},
			{
				Name:      "mysql-mysql1",
				Image:     "test/mysqlimage",
				Resources: resources,
				Env: []corev1.EnvVar{
					{
						Name: EnvMySQLHostname,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-host",
							},
						},
					},
					{
						Name: EnvMySQLDatabase,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-db",
							},
						},
					},
					{
						Name: EnvMySQLPort,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-port",
							},
						},
					},
					{
						Name: EnvMySQLUsername,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-user",
							},
						},
					},
					{
						Name: EnvMySQLPassword,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-pass",
							},
						},
					},
				},
				WorkingDir: "/home/test",
				Command: []string{
					"database-backup",
				},
				Args: []string{
					"mysql/mysql1.sql",
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      VolumeMySQL,
						MountPath: "/home/test/mysql",
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:       ResticBackupContainerName,
				Image:      "test/image",
				Resources:  resources,
				WorkingDir: "/home/test",
				Command: []string{
					"/bin/sh", "-c",
				},
				Args: []string{
					"restic --verbose --tag=tag1 backup . --exclude volume/*/*/php --exclude volume/*/*/css --exclude volume/*/*/js",
				},
				Env: []corev1.EnvVar{
					{
						Name:  EnvResticRepository,
						Value: "/srv/backups/test-namespace/test-site-id",
					},
					{
						Name:  EnvResticPasswordFile,
						Value: "/etc/restic/password",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      VolumeMySQL,
						MountPath: "/home/test/mysql",
					},
					{
						Name:      "volume-volume1",
						MountPath: "/home/test/volume/volume1",
						ReadOnly:  true,
					},
					{
						Name:      VolumeSecrets,
						MountPath: SecretDir,
						ReadOnly:  true,
					},
					{
						Name:      VolumeRepository,
						MountPath: ResticRepoDir,
					},
				},
			},
		},
		Volumes: []corev1.Volume{
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
			{
				Name: VolumeSecrets,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: &mode,
						SecretName:  ResticSecretPasswordName,
					},
				},
			},
			{
				Name: "volume-volume1",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "claim-volume1",
					},
				},
			},
		},
	}

	spec, _ := PodSpecBackup(&backup, params, "test-site-id")
	assert.Equal(t, expected, spec)
}

func TestPodSpecRestore(t *testing.T) {
	var params = PodSpecParams{
		CPU:         "100m",
		Memory:      "512Mi",
		ResticImage: "test/image",
		MySQLImage:  "test/mysqlimage",
		WorkingDir:  "/home/test",
		Tags:        []string{"tag1"},
	}
	restore := extensionv1.Restore{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-restore",
			Namespace: "test-namespace",
		},
		Spec: extensionv1.RestoreSpec{
			BackupName: "test-backup",
			Volumes: map[string]metav1_shepherd.SpecVolume{
				"volume1": {
					ClaimName: "claim-volume1",
				},
			},
			MySQL: map[string]metav1_shepherd.SpecMySQL{
				"mysql1": {
					Secret: metav1_shepherd.SpecMySQLSecret{
						Name: "secret1",
						Keys: metav1_shepherd.SpecMySQLSecretKeys{
							Username: "mysql-user",
							Password: "mysql-pass",
							Database: "mysql-db",
							Hostname: "mysql-host",
							Port:     "mysql-port",
						},
					},
				},
			},
		},
	}

	dc := osv1.DeploymentConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-dc",
			Namespace: "test-namespace",
		},
		Spec: osv1.DeploymentConfigSpec{
			Template: &corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test-dc",
							Image: "test/deploy-image",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "different-named-volume",
									MountPath: "/testmount",
								},
								{
									Name:      "another-unrelated-volume",
									MountPath: "/testmount2",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "foo",
									Value: "bar",
								},
								{
									Name:  "baz",
									Value: "blop",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							// Volume with the same claim name as the restore to test names are overwritten.
							Name: "different-named-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "claim-volume1",
								},
							},
						},
						{
							Name: "another-unrelated-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
	cpu, _ := resource.ParseQuantity(params.CPU)
	memory, _ := resource.ParseQuantity(params.Memory)
	mode := corev1.ConfigMapVolumeSourceDefaultMode
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
	expected := corev1.PodSpec{
		RestartPolicy: corev1.RestartPolicyNever,
		InitContainers: []corev1.Container{
			{
				Name:       "restic-restore-mysql1",
				Image:      "test/image",
				Resources:  resources,
				WorkingDir: "/home/test",
				Command: []string{
					"/bin/sh", "-c",
				},
				Args: []string{
					"restic dump abcd1234 /mysql/mysql1.sql > ./mysql/mysql1.sql",
				},
				Env: []corev1.EnvVar{
					{
						Name:  EnvResticRepository,
						Value: "/srv/backups/test-namespace/test-site-id",
					},
					{
						Name:  EnvResticPasswordFile,
						Value: "/etc/restic/password",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      VolumeMySQL,
						MountPath: "/home/test/mysql",
					},
					{
						Name:      VolumeSecrets,
						MountPath: SecretDir,
						ReadOnly:  true,
					},
					{
						Name:      VolumeRepository,
						MountPath: ResticRepoDir,
					},
				},
			},
			{
				Name:       "restic-import-mysql1",
				Image:      "test/mysqlimage",
				Resources:  resources,
				WorkingDir: "/home/test",
				Command: []string{
					"database-restore",
				},
				Args: []string{
					"mysql/mysql1.sql",
				},
				Env: []corev1.EnvVar{
					{
						Name: EnvMySQLHostname,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-host",
							},
						},
					},
					{
						Name: EnvMySQLDatabase,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-db",
							},
						},
					},
					{
						Name: EnvMySQLPort,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-port",
							},
						},
					},
					{
						Name: EnvMySQLUsername,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-user",
							},
						},
					},
					{
						Name: EnvMySQLPassword,
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "secret1",
								},
								Key: "mysql-pass",
							},
						},
					},
					{
						Name:  EnvResticRepository,
						Value: "/srv/backups/test-namespace/test-site-id",
					},
					{
						Name:  EnvResticPasswordFile,
						Value: "/etc/restic/password",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      VolumeMySQL,
						MountPath: "/home/test/mysql",
					},
					{
						Name:      VolumeSecrets,
						MountPath: SecretDir,
						ReadOnly:  true,
					},
					{
						Name:      VolumeRepository,
						MountPath: ResticRepoDir,
					},
				},
			},
			{
				Name:       "restic-restore-volumes",
				Image:      "test/image",
				Resources:  resources,
				WorkingDir: "/home/test",
				Command: []string{
					"/bin/sh", "-c",
				},
				Args: []string{
					"restic restore abcd1234 --target . --include /volume/volume1",
				},
				Env: []corev1.EnvVar{
					{
						Name:  EnvResticRepository,
						Value: "/srv/backups/test-namespace/test-site-id",
					},
					{
						Name:  EnvResticPasswordFile,
						Value: "/etc/restic/password",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "volume-volume1",
						MountPath: "/home/test/volume/volume1",
					},
					{
						Name:      VolumeSecrets,
						MountPath: SecretDir,
						ReadOnly:  true,
					},
					{
						Name:      VolumeRepository,
						MountPath: ResticRepoDir,
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:       "restore-deploy",
				Image:      "test/deploy-image",
				Resources:  resources,
				WorkingDir: WebDirectory,
				Command: []string{
					"/bin/sh", "-c",
				},
				Args: []string{
					"drush -r /code/web cr && drush -r /code/web -y updb && robo config:import-plus && drush -r /code/web cr",
				},
				Env: []corev1.EnvVar{
					{
						Name:  "foo",
						Value: "bar",
					},
					{
						Name:  "baz",
						Value: "blop",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "volume-volume1",
						MountPath: "/testmount",
					},
					{
						Name:      "another-unrelated-volume",
						MountPath: "/testmount2",
					},
				},
			},
		},
		Volumes: []corev1.Volume{
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
			{
				Name: VolumeSecrets,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: &mode,
						SecretName:  ResticSecretPasswordName,
					},
				},
			},
			{
				Name: "volume-volume1",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "claim-volume1",
					},
				},
			},
			{
				Name: "another-unrelated-volume",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		},
	}

	spec, _ := PodSpecRestore(&restore, &dc, "abcd1234", params, "test-site-id")
	assert.Equal(t, expected, spec)
}
