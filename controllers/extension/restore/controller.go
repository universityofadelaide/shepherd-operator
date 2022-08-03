/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	osv1 "github.com/openshift/api/apps/v1"
	osv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpdmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	"github.com/universityofadelaide/shepherd-operator/internal/events"
	"github.com/universityofadelaide/shepherd-operator/internal/helper"
	podutils "github.com/universityofadelaide/shepherd-operator/internal/k8s/pod"
)

const (
	// ControllerName is used to identify this controller in logs and events.
	ControllerName = "restore-controller"

	// EnvAWSAccessKeyID for authentication.
	EnvAWSAccessKeyID = "AWS_ACCESS_KEY_ID"
	// EnvAWSSecretAccessKey for authentication.
	EnvAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	// EnvAWSRegion for authentication.
	EnvAWSRegion = "AWS_DEFAULT_REGION"

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
	VolumeMySQL = "mysql"

	// WebDirectory is working directory for the restore deployment step.
	WebDirectory = "/code"
)

// Reconciler reconciles a Restore object
type Reconciler struct {
	client.Client
	OpenShift osv1client.AppsV1Interface
	Config    *rest.Config
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Params    Params
}

// Params used by this controller.
type Params struct {
	ResourceRequirements corev1.ResourceRequirements
	WorkingDir           string
	// MySQL params used by this controller.
	MySQL MySQL
	// AWS params used by this controller.
	AWS AWS
}

// MySQL params used by this controller.
type MySQL struct {
	Image string
}

// AWS params used by this controller.
type AWS struct {
	BucketName     string
	Image          string
	FieldKeyID     string
	FieldAccessKey string
	Region         string
}

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=restores/finalizers,verbs=update

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconcile loop")

	restore := &extensionv1.Restore{}

	err := r.Get(ctx, req.NamespacedName, restore)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	backup := &extensionv1.Backup{}

	err = r.Get(ctx, types.NamespacedName{
		Name:      restore.Spec.BackupName,
		Namespace: restore.ObjectMeta.Namespace,
	}, backup)
	if err != nil {
		if kerrors.IsNotFound(err) {
			r.Recorder.Eventf(restore, corev1.EventTypeNormal, events.EventError, "Backup not found: %s", restore.Spec.BackupName)
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	switch backup.Status.Phase {
	case shpdmetav1.PhaseFailed:
		logger.Info(fmt.Sprintf("Skipping restore %s because the backup %s failed", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return reconcile.Result{}, nil
	case shpdmetav1.PhaseNew:
		// Requeue the operation for 30 seconds if the backup is new.
		logger.Info(fmt.Sprintf("Requeueing restore %s because the backup %s is New", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 30}, nil
	case shpdmetav1.PhaseInProgress:
		// Requeue the operation for 15 seconds if the backup is still in progress.
		logger.Info(fmt.Sprintf("Requeueing restore %s because the backup %s is In Progress", restore.ObjectMeta.Name, backup.ObjectMeta.Name))
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 15}, nil

	}

	// Catch-all for any other non Completed phases.
	if backup.Status.Phase != shpdmetav1.PhaseCompleted {
		logger.Info(fmt.Sprintf("Skipping restore %s because the backup %s is in an unknown state: %s", restore.ObjectMeta.Name, backup.ObjectMeta.Name, backup.Status.Phase))
		return reconcile.Result{}, nil
	}

	if _, found := restore.ObjectMeta.GetLabels()["site"]; !found {
		// @todo add some info to the status identifying the restore failed
		logger.Info(fmt.Sprintf("Restore %s doesn't have a site label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}
	// TODO: Add environment to spec so we don't have to derive the deploymentconfig name.
	if _, found := restore.ObjectMeta.GetLabels()["environment"]; !found {
		logger.Info(fmt.Sprintf("Restore %s doesn't have a environment label, skipping.", restore.ObjectMeta.Name))
		return reconcile.Result{}, nil
	}

	dcName := fmt.Sprintf("node-%s", restore.ObjectMeta.GetLabels()["environment"])

	dc, err := r.OpenShift.DeploymentConfigs(restore.ObjectMeta.Namespace).Get(ctx, dcName, metav1.GetOptions{})
	if err != nil {
		// Don't throw an error here to account for restores that were ted before an environment was deleted.
		return reconcile.Result{}, nil
	}

	secret, err := r.createSecret(ctx, backup, r.Params.AWS.FieldKeyID, r.Params.AWS.FieldAccessKey)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create Secret: %w", err)
	}

	status, err := r.createPod(ctx, restore, secret, dc)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create Pod: %w", err)
	}

	err = r.updateStatus(ctx, logger, restore, status)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update Backup status: %w", err)
	}

	logger.Info("Reconcile finished")

	return reconcile.Result{}, nil
}

// Creates Secret object based on the provided Spec configuration.
func (r *Reconciler) createSecret(ctx context.Context, backup *extensionv1.Backup, key, access string) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("restore-%s", backup.ObjectMeta.Name),
			Namespace: backup.ObjectMeta.Namespace,
		},
		Data: map[string][]byte{
			EnvAWSAccessKeyID:     []byte(key),
			EnvAWSSecretAccessKey: []byte(access),
		},
	}

	if err := controllerutil.SetControllerReference(backup, secret, r.Scheme); err != nil {
		return nil, err
	}

	if err := r.Create(ctx, secret); client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	return secret, nil
}

// Creates Pod objects based on the provided Spec configuration.
func (r *Reconciler) createPod(ctx context.Context, restore *extensionv1.Restore, secret *corev1.Secret, dc *osv1.DeploymentConfig) (extensionv1.RestoreStatus, error) {
	var initContainers []corev1.Container
	var containers []corev1.Container

	// InitContainer which restores db to emptydir volume.
	for mysqlName, mysqlStatus := range restore.Spec.MySQL {
		initContainers = append(initContainers, corev1.Container{
			Name:       fmt.Sprintf("restore-%s", mysqlName),
			Image:      r.Params.AWS.Image,
			Resources:  r.Params.ResourceRequirements,
			WorkingDir: r.Params.WorkingDir,
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				// @todo, Remove hardcoded command.
				// @todo, Determine requirements for S3 path.
				helper.TprintfMustParse(
					"aws s3 cp s3://{{.BucketName}}/{{.Namespace}}/{{.BackupName}}/mysql/{{.MySQLName}}.sql mysql/{{.MySQLName}}.sql",
					map[string]interface{}{
						"BucketName": r.Params.AWS.BucketName,
						"Namespace":  restore.ObjectMeta.Namespace,
						"BackupName": restore.Spec.BackupName,
						"MySQLName":  mysqlName,
					},
				),
			},
			Env: []corev1.EnvVar{
				{
					Name: EnvAWSAccessKeyID,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: secret.ObjectMeta.Name,
							},
							Key: EnvAWSAccessKeyID,
						},
					},
				},
				{
					Name: EnvAWSSecretAccessKey,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: secret.ObjectMeta.Name,
							},
							Key: EnvAWSSecretAccessKey,
						},
					},
				},
				{
					Name:  EnvAWSRegion,
					Value: r.Params.AWS.Region,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", r.Params.WorkingDir),
				},
			},
		})

		initContainers = append(initContainers, corev1.Container{
			Name:       fmt.Sprintf("import-%s", mysqlName),
			Image:      r.Params.MySQL.Image,
			Resources:  r.Params.ResourceRequirements,
			WorkingDir: r.Params.WorkingDir,
			Command: []string{
				"database-restore",
			},
			Args: []string{
				fmt.Sprintf("mysql/%s.sql", mysqlName),
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      VolumeMySQL,
					MountPath: fmt.Sprintf("%s/mysql", r.Params.WorkingDir),
				},
			},
			Env: []corev1.EnvVar{
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
			},
		})
	}

	// Volume definitions for the pod.
	specVolumes := []corev1.Volume{
		{
			Name: VolumeMySQL,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

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

		// Container which restores volumes.
		initContainers = append(initContainers, corev1.Container{
			Name:       "restore-volumes",
			Image:      r.Params.AWS.Image,
			Resources:  r.Params.ResourceRequirements,
			WorkingDir: r.Params.WorkingDir,
			Command: []string{
				"/bin/sh", "-c",
			},
			Args: []string{
				// @todo, Remove hardcoded command.
				// @todo, Determine requirements for S3 path.
				helper.TprintfMustParse(
					"aws s3 sync --delete s3://{{.BucketName}}/{{.Namespace}}/{{.BackupName}}/{{.VolumeName}}/ {{.MountPath}}/",
					map[string]interface{}{
						"BucketName": r.Params.AWS.BucketName,
						"Namespace":  restore.ObjectMeta.Namespace,
						"BackupName": restore.Spec.BackupName,
						"VolumeName": volumeName,
						"MountPath":  fmt.Sprintf("%s/volume/%s", r.Params.WorkingDir, volumeName),
					},
				),
			},
			Env: []corev1.EnvVar{
				{
					Name: EnvAWSAccessKeyID,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: secret.ObjectMeta.Name,
							},
							Key: EnvAWSAccessKeyID,
						},
					},
				},
				{
					Name: EnvAWSSecretAccessKey,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: secret.ObjectMeta.Name,
							},
							Key: EnvAWSSecretAccessKey,
						},
					},
				},
				{
					Name:  EnvAWSRegion,
					Value: r.Params.AWS.Region,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      fmt.Sprintf("volume-%s", volumeName),
					MountPath: fmt.Sprintf("%s/volume/%s", r.Params.WorkingDir, volumeName),
					ReadOnly:  false,
				},
			},
		})
	}

	var status extensionv1.RestoreStatus

	dcContainer, err := getWebContainerFromDc(dc)
	if err != nil {
		return status, err
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
	// @todo, Try and make this into a reusable CRD.
	containers = append(containers, corev1.Container{
		Name:       "restore-deploy",
		Image:      dcContainer.Image,
		Resources:  r.Params.ResourceRequirements,
		WorkingDir: WebDirectory,
		Command: []string{
			"/bin/sh", "-c",
		},
		Args: []string{
			helper.TprintfMustParse(
				"drush -r {{.WebDir}}/web cr && drush -r {{.WebDir}}/web -y updb && robo config:import-plus && drush -r {{.WebDir}}/web cr",
				map[string]interface{}{
					"WebDir": WebDirectory,
				},
			),
		},
		Env:          dcContainer.Env,
		VolumeMounts: dcVolumeMounts,
	})

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("restore-%s", restore.ObjectMeta.Name),
			Namespace: restore.ObjectMeta.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:  corev1.RestartPolicyNever,
			InitContainers: initContainers,
			Containers:     containers,
			Volumes:        specVolumes,
		},
	}

	if err := controllerutil.SetControllerReference(restore, pod, r.Scheme); err != nil {
		return status, err
	}

	if err := r.Create(ctx, pod); client.IgnoreNotFound(err) != nil {
		return status, err
	}

	if err := r.Get(ctx, types.NamespacedName{
		Namespace: pod.ObjectMeta.Namespace,
		Name:      pod.ObjectMeta.Name,
	}, pod); err != nil {
		return status, err
	}

	status.Phase = podutils.GetPhase(pod.Status)
	status.StartTime = pod.Status.StartTime
	status.CompletionTime = podutils.CompletionTime(pod)

	return status, nil
}

// Update the Backup status.
func (r *Reconciler) updateStatus(ctx context.Context, log logr.Logger, restore *extensionv1.Restore, status extensionv1.RestoreStatus) error {
	diff := deep.Equal(restore.Status, status)
	if diff == nil {
		return nil
	}

	log.Info(fmt.Sprintf("Status change dectected: %s", diff))

	restore.Status = status

	return r.Status().Update(ctx, restore)
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

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Restore{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
