package drupal

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
	"github.com/skpr/operator/pkg/utils/k8s/generate"
)

const (
	// VolumePermissionsContainer identifier for the permissions CronJob container.
	VolumePermissionsContainer = "permissions"
	// VolumePermissionsDefaultDeadline sets deadline if not specified.
	VolumePermissionsDefaultDeadline = 600
	// VolumePermissionsDefaultKeepSuccess sets KeepSuccess if not specified.
	VolumePermissionsDefaultKeepSuccess = 1
	// VolumePermissionsDefaultKeepFailed sets KeepFailed if not specified.
	VolumePermissionsDefaultKeepFailed = 1
	// VolumePermissionsDefaultRetries sets Retries if not specified.
	VolumePermissionsDefaultRetries = 2
)

// Helper function to create CronJobs to ensure permissions for a Volume.
func buildVolumePermissions(name string, drupal *appv1beta1.Drupal, pvc *corev1.PersistentVolumeClaim, volume appv1beta1.DrupalSpecVolume) (*batchv1beta1.CronJob, *batchv1beta1.CronJob, *batchv1beta1.CronJob, error) {
	labels := map[string]string{
		LabelAppName:  drupal.ObjectMeta.Name,
		LabelAppType:  Application,
		LabelAppLayer: LayerVolume,
	}

	if volume.Permissions.CronJob.Deadline == 0 {
		volume.Permissions.CronJob.Deadline = VolumePermissionsDefaultDeadline
	}

	if volume.Permissions.CronJob.Retries == 0 {
		volume.Permissions.CronJob.Retries = VolumePermissionsDefaultRetries
	}

	if volume.Permissions.CronJob.KeepSuccess == 0 {
		volume.Permissions.CronJob.KeepSuccess = VolumePermissionsDefaultKeepSuccess
	}

	if volume.Permissions.CronJob.KeepFailed == 0 {
		volume.Permissions.CronJob.KeepFailed = VolumePermissionsDefaultKeepFailed
	}

	grace := int64(corev1.DefaultTerminationGracePeriodSeconds)

	files := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-perm-files", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:                   volume.Permissions.CronJob.Schedule,
			StartingDeadlineSeconds:    &volume.Permissions.CronJob.Deadline,
			ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
			SuccessfulJobsHistoryLimit: &volume.Permissions.CronJob.KeepSuccess,
			FailedJobsHistoryLimit:     &volume.Permissions.CronJob.KeepFailed,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: drupal.ObjectMeta.Namespace,
					Labels:    labels,
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &volume.Permissions.CronJob.Retries,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: drupal.ObjectMeta.Namespace,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            VolumePermissionsContainer,
									Image:           volume.Permissions.CronJob.Image,
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										ReadOnlyRootFilesystem: &volume.Permissions.CronJob.ReadOnly,
									},
									WorkingDir: volume.Path,
									Command: []string{
										"/bin/sh", "-c",
									},
									Args: []string{
										fmt.Sprintf("find . -type f -exec chmod %d {} +", volume.Permissions.File),
									},
									VolumeMounts: []corev1.VolumeMount{
										generate.Mount(VolumePublic, volume.Path, false),
									},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								},
							},
							Volumes: []corev1.Volume{
								generate.VolumeClaim(VolumePublic, pvc.ObjectMeta.Name),
							},
							// The below are fields which need to be set so we can perform an "deep equal"
							// without always having difference.
							SecurityContext:               &corev1.PodSecurityContext{},
							SchedulerName:                 corev1.DefaultSchedulerName,
							DNSPolicy:                     corev1.DNSClusterFirst,
							TerminationGracePeriodSeconds: &grace,
							RestartPolicy:                 corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	dirs := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-perm-dirs", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:                   volume.Permissions.CronJob.Schedule,
			StartingDeadlineSeconds:    &volume.Permissions.CronJob.Deadline,
			ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
			SuccessfulJobsHistoryLimit: &volume.Permissions.CronJob.KeepSuccess,
			FailedJobsHistoryLimit:     &volume.Permissions.CronJob.KeepFailed,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: drupal.ObjectMeta.Namespace,
					Labels:    labels,
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &volume.Permissions.CronJob.Retries,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: drupal.ObjectMeta.Namespace,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            VolumePermissionsContainer,
									Image:           volume.Permissions.CronJob.Image,
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										ReadOnlyRootFilesystem: &volume.Permissions.CronJob.ReadOnly,
									},
									WorkingDir: volume.Path,
									Command: []string{
										"/bin/sh", "-c",
									},
									Args: []string{
										fmt.Sprintf("find . -type d -exec chmod %d {} +", volume.Permissions.Directory),
									},
									VolumeMounts: []corev1.VolumeMount{
										generate.Mount(VolumePublic, volume.Path, false),
									},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								},
							},
							Volumes: []corev1.Volume{
								generate.VolumeClaim(VolumePublic, pvc.ObjectMeta.Name),
							},
							// The below are fields which need to be set so we can perform an "deep equal"
							// without always having difference.
							SecurityContext:               &corev1.PodSecurityContext{},
							SchedulerName:                 corev1.DefaultSchedulerName,
							DNSPolicy:                     corev1.DNSClusterFirst,
							TerminationGracePeriodSeconds: &grace,
							RestartPolicy:                 corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	owner := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-perm-owner", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:                   volume.Permissions.CronJob.Schedule,
			StartingDeadlineSeconds:    &volume.Permissions.CronJob.Deadline,
			ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
			SuccessfulJobsHistoryLimit: &volume.Permissions.CronJob.KeepSuccess,
			FailedJobsHistoryLimit:     &volume.Permissions.CronJob.KeepFailed,
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: drupal.ObjectMeta.Namespace,
					Labels:    labels,
				},
				Spec: batchv1.JobSpec{
					BackoffLimit: &volume.Permissions.CronJob.Retries,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: drupal.ObjectMeta.Namespace,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            VolumePermissionsContainer,
									Image:           volume.Permissions.CronJob.Image,
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										ReadOnlyRootFilesystem: &volume.Permissions.CronJob.ReadOnly,
									},
									WorkingDir: volume.Path,
									Command: []string{
										"/bin/sh", "-c",
									},
									Args: []string{
										fmt.Sprintf("chown -R %s:%s .", volume.Permissions.User, volume.Permissions.Group),
									},
									VolumeMounts: []corev1.VolumeMount{
										generate.Mount(VolumePublic, volume.Path, false),
									},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								},
							},
							Volumes: []corev1.Volume{
								generate.VolumeClaim(VolumePublic, pvc.ObjectMeta.Name),
							},
							// The below are fields which need to be set so we can perform an "deep equal"
							// without always having difference.
							SecurityContext:               &corev1.PodSecurityContext{},
							SchedulerName:                 corev1.DefaultSchedulerName,
							DNSPolicy:                     corev1.DNSClusterFirst,
							TerminationGracePeriodSeconds: &grace,
							RestartPolicy:                 corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	return files, dirs, owner, nil
}
