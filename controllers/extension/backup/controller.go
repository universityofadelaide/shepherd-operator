package backup

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	extensionv1 "github.com/universityofadelaide/shepherd-operator/apis/extension/v1"
	shpdmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	awscli "github.com/universityofadelaide/shepherd-operator/internal/aws/cli"
	podutils "github.com/universityofadelaide/shepherd-operator/internal/k8s/pod"
)

const (
	// ControllerName is used to identify this controller in logs and events.
	ControllerName = "backup-controller"

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
)

// Reconciler reconciles a Backup object
type Reconciler struct {
	client.Client
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	Params   Params
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
	Endpoint       string
	BucketName     string
	Image          string
	FieldKeyID     string
	FieldAccessKey string
	Region         string
}

//+kubebuilder:rbac:groups=batch,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=pods/status,verbs=get
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=extension.shepherd,resources=backups/finalizers,verbs=update

// Reconcile a Backup object
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconcile loop")

	backup := &extensionv1.Backup{}

	if err := r.Get(ctx, req.NamespacedName, backup); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Backup has completed or failed, return early.
	if backup.Status.Phase == shpdmetav1.PhaseCompleted || backup.Status.Phase == shpdmetav1.PhaseFailed {
		return reconcile.Result{}, nil
	}

	if backup.Spec.Type == "" {
		backup.Spec.Type = extensionv1.BackupTypeDefault
	}

	err := r.createSecret(ctx, backup, r.Params.AWS.FieldKeyID, r.Params.AWS.FieldAccessKey)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create Secret: %w", err)
	}

	status, err := r.createPod(ctx, backup)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create Pod: %w", err)
	}

	err = r.updateStatus(ctx, logger, backup, status)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update Backup status: %w", err)
	}

	logger.Info("Finished reconcile loop")

	return ctrl.Result{}, nil
}

// Creates Secret object based on the provided Spec configuration.
func (r *Reconciler) createSecret(ctx context.Context, backup *extensionv1.Backup, key, access string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(backup),
			Namespace: backup.ObjectMeta.Namespace,
		},
		Data: map[string][]byte{
			EnvAWSAccessKeyID:     []byte(key),
			EnvAWSSecretAccessKey: []byte(access),
		},
	}

	if err := controllerutil.SetControllerReference(backup, secret, r.Scheme); err != nil {
		return err
	}

	err := r.Create(ctx, secret)

	if kerrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

// Creates Pod objects based on the provided Spec configuration.
func (r *Reconciler) createPod(ctx context.Context, backup *extensionv1.Backup) (extensionv1.BackupStatus, error) {
	cmd := awscli.CommandParams{
		Endpoint:  r.Params.AWS.Endpoint,
		Service:   "s3",
		Operation: "sync",
		Args: []string{
			".", fmt.Sprintf("s3://%s/%s/%s/%s", r.Params.AWS.BucketName, backup.Spec.Type, backup.ObjectMeta.Namespace, backup.ObjectMeta.Name),
		},
	}

	// @todo, This should be configured at the object level.
	exclude := []string{
		"volume/*/*/php",
		"volume/*/*/css",
		"volume/*/*/js",
	}

	for _, exclude := range exclude {
		cmd.Args = append(cmd.Args, "--exclude", exclude)
	}

	// Container responsible for uploading database and files to AWS S3.
	upload := corev1.Container{
		Name:            "aws-s3-sync",
		Image:           r.Params.AWS.Image,
		ImagePullPolicy: corev1.PullAlways,
		Resources:       r.Params.ResourceRequirements,
		WorkingDir:      r.Params.WorkingDir,
		Args:            awscli.Command(cmd),
		Env: []corev1.EnvVar{
			{
				Name: EnvAWSAccessKeyID,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: getName(backup),
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
							Name: getName(backup),
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
	}

	for volumeName := range backup.Spec.Volumes {
		upload.VolumeMounts = append(upload.VolumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("volume-%s", volumeName),
			MountPath: fmt.Sprintf("%s/volume/%s", r.Params.WorkingDir, volumeName),
			ReadOnly:  true,
		})
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(backup),
			Namespace: backup.ObjectMeta.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				upload,
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
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	for volumeName, volumeSpec := range backup.Spec.Volumes {
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: fmt.Sprintf("volume-%s", volumeName),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeSpec.ClaimName,
				},
			},
		})
	}

	for mysqlName, mysqlStatus := range backup.Spec.MySQL {
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:       fmt.Sprintf("mysql-%s", mysqlName),
			Image:      r.Params.MySQL.Image,
			Resources:  r.Params.ResourceRequirements,
			WorkingDir: r.Params.WorkingDir,
			Command: []string{
				"bash",
				"-c",
			},
			Args: []string{
				// @todo, Remove hardcoded command and path.
				fmt.Sprintf("database-backup > mysql/%s.sql", mysqlName),
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
			VolumeMounts: []corev1.VolumeMount{
				{
					Name: VolumeMySQL,
					// @todo, Remove hardcoded mysql path.
					MountPath: fmt.Sprintf("%s/mysql", r.Params.WorkingDir),
				},
			},
		})
	}

	var status extensionv1.BackupStatus

	if err := controllerutil.SetControllerReference(backup, pod, r.Scheme); err != nil {
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
func (r *Reconciler) updateStatus(ctx context.Context, log logr.Logger, backup *extensionv1.Backup, status extensionv1.BackupStatus) error {
	diff := deep.Equal(backup.Status, status)
	if diff == nil {
		return nil
	}

	log.Info(fmt.Sprintf("Status change dectected: %s", diff))

	backup.Status = status

	return r.Status().Update(ctx, backup)
}

// Helper function to get a resource name.
func getName(backup *extensionv1.Backup) string {
	return fmt.Sprintf("backup-%s", backup.ObjectMeta.Name)
}

// SetupWithManager will setup the controller.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionv1.Backup{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
