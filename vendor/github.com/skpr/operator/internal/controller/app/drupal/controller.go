package drupal

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appv1beta1 "github.com/skpr/operator/pkg/apis/app/v1beta1"
	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
	"github.com/skpr/operator/pkg/mysql"
	"github.com/skpr/operator/pkg/utils/controller/logger"
	deploymentutils "github.com/skpr/operator/pkg/utils/k8s/deployment"
	"github.com/skpr/operator/pkg/utils/k8s/generate"
	k8ssync "github.com/skpr/operator/pkg/utils/k8s/sync"
	"github.com/skpr/operator/pkg/utils/prometheus"
)

const (
	// ControllerName used for identifying which controller is performing an operation.
	ControllerName = "drupal-controller"
)

const (
	// Application identifier.
	Application = "drupal"
)

const (
	// LabelAppName for discovery.
	LabelAppName = "app_name"
	// LabelAppType for discovery.
	LabelAppType = "app_type"
	// LabelAppLayer for discovery.
	LabelAppLayer = "app_layer"
)

const (
	// LayerNginx for identifying Nginx objects.
	LayerNginx = "nginx"
	// LayerFPM for identifying FPM objects.
	LayerFPM = "fpm"
	// LayerBatch for identifying Batch objects.
	LayerBatch = "batch"
	// LayerVolume for identifying Volume objects.
	LayerVolume = "volume"
)

const (
	// PodEnvSkipperEnv identifies which environment an application is running.
	PodEnvSkipperEnv = "SKIPPER_ENV"
	// PodEnvNewRelicApp is a required environment variable for New Relic monitoring.
	PodEnvNewRelicApp = "NEW_RELIC_APP_NAME"
	// PodEnvNewRelicLicense is a required environment variable for New Relic monitoring.
	PodEnvNewRelicLicense = "NEW_RELIC_LICENSE_KEY"
	// PodEnvNewRelicEnabled is a required environment variable for New Relic monitoring.
	PodEnvNewRelicEnabled = "NEW_RELIC_ENABLED"
)

const (
	// PodContainerNginx for identifying the Nginx container in a pod.
	PodContainerNginx = "nginx"
	// PodContainerFPM for identifying the FPM container in a pod.
	PodContainerFPM = "fpm"
	// PodContainerCLI for identifying the CLI container in a pod.
	PodContainerCLI = "cli"
	// PodContainerExporter for identifying the Exporter container in a pod.
	PodContainerExporter = "exporter"
)

const (
	// VolumePublic identifier.
	VolumePublic = "public"
	// VolumePrivate identifier.
	VolumePrivate = "private"
	// VolumeTemporary identifier.
	VolumeTemporary = "temporary"
	// VolumeConfigDefault identifier.
	VolumeConfigDefault = "config-default"
	// VolumeConfigOverride identifier.
	VolumeConfigOverride = "config-override"
	// VolumeSecretDefault identifier.
	VolumeSecretDefault = "secret-default"
	// VolumeSecretOverride identifier.
	VolumeSecretOverride = "secret-override"
)

const (
	// ConfigMapMountPublic declares which config Drupal can use to store public files.
	ConfigMapMountPublic = "mount.public"
	// ConfigMapMountPrivate declares which config Drupal can use to store private files.
	ConfigMapMountPrivate = "mount.private"
	// ConfigMapMountTemporary declares which config Drupal can use to store temporary files.
	ConfigMapMountTemporary = "mount.temporary"
)

const (
	// SecretPrometheusToken identifier for applications.
	SecretPrometheusToken = "prometheus.token"
)

// Add creates a new Drupal Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, exporters Exporters) error {
	return add(mgr, newReconciler(mgr, exporters))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, exporters Exporters) reconcile.Reconciler {
	return &ReconcileDrupal{
		Client:    mgr.GetClient(),
		recorder:  mgr.GetRecorder(ControllerName),
		scheme:    mgr.GetScheme(),
		exporters: exporters,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch Deployment changes.
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1beta1.Drupal{},
	})
	if err != nil {
		return err
	}

	// Watch Drupal changes.
	return c.Watch(&source.Kind{Type: &appv1beta1.Drupal{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileDrupal{}

// ReconcileDrupal reconciles a Drupal object
type ReconcileDrupal struct {
	client.Client
	recorder  record.EventRecorder
	scheme    *runtime.Scheme
	exporters Exporters
}

// Exporters which are used for Drupal metrics.
type Exporters struct {
	Nginx Exporter
	FPM   Exporter
}

// Exporter which is used for Drupal metrics.
type Exporter struct {
	Image  string
	Port   string
	CPU    resource.Quantity
	Memory resource.Quantity
}

// Reconcile reads that state of the cluster for a Drupal object and makes changes based on the state read
// and what is in the Drupal.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.skpr.io,resources=drupals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.skpr.io,resources=drupals/status,verbs=get;update;patch
func (r *ReconcileDrupal) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := logger.New(ControllerName, request.Namespace, request.Name)

	log.Info("Starting Reconcile Loop")

	drupal := &appv1beta1.Drupal{}
	err := r.Get(context.TODO(), request.NamespacedName, drupal)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	log.Info("Syncing Drupal Deployment")

	status, err := r.Sync(log, drupal)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "sync status failed")
	}

	err = r.SyncStatus(drupal, status)
	if err != nil {
		log.Error(err, "Status status failed")
		return reconcile.Result{}, errors.Wrap(err, "sync status failed")
	}

	log.Info("Reconcile Loop Finished")

	return reconcile.Result{RequeueAfter: time.Second * 15}, err
}

// Sync all Kubernetes objects and return the status of the Drupal deployment.
func (r *ReconcileDrupal) Sync(log log.Logger, drupal *appv1beta1.Drupal) (appv1beta1.DrupalStatus, error) {
	var status appv1beta1.DrupalStatus

	// Common base name for all resources.
	name := fmt.Sprintf("%s-%s", Application, drupal.ObjectMeta.Name)

	commonlabels := map[string]string{
		LabelAppName: drupal.ObjectMeta.Name,
		LabelAppType: Application,
	}

	configMapDefault := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-default", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Data: map[string]string{
			ConfigMapMountPublic:    drupal.Spec.Volume.Public.Path,
			ConfigMapMountPrivate:   drupal.Spec.Volume.Private.Path,
			ConfigMapMountTemporary: drupal.Spec.Volume.Temporary.Path,
		},
	}

	secretDefault := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-default", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Data: make(map[string][]byte),
	}

	if drupal.Spec.Prometheus.Token != "" {
		secretDefault.Data[SecretPrometheusToken] = []byte(drupal.Spec.Prometheus.Token)
	}

	pvcPublic := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, VolumePublic),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(drupal.Spec.Volume.Public.Amount),
				},
			},
			StorageClassName: &drupal.Spec.Volume.Public.Class,
		},
	}

	pvcPrivate := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, VolumePrivate),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(drupal.Spec.Volume.Private.Amount),
				},
			},
			StorageClassName: &drupal.Spec.Volume.Private.Class,
		},
	}

	pvcTemporary := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, VolumeTemporary),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(drupal.Spec.Volume.Temporary.Amount),
				},
			},
			StorageClassName: &drupal.Spec.Volume.Temporary.Class,
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPublic, k8ssync.PersistentVolumeClaim(drupal, pvcPublic.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync PersistentVolumeClaim")
	}
	log.Infof("Synced PersistentVolumeClaim %s with status: %s", pvcPublic.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPrivate, k8ssync.PersistentVolumeClaim(drupal, pvcPrivate.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync PersistentVolumeClaim")
	}
	log.Infof("Synced PersistentVolumeClaim %s with status: %s", pvcPrivate.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcTemporary, k8ssync.PersistentVolumeClaim(drupal, pvcTemporary.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync PersistentVolumeClaim")
	}
	log.Infof("Synced PersistentVolumeClaim %s with status: %s", pvcTemporary.ObjectMeta.Name, result)

	pvcPublicFiles, pvcPublicDirs, pvcPublicOwner, err := buildVolumePermissions(fmt.Sprintf("%s-public", name), drupal, pvcPublic, drupal.Spec.Volume.Public)
	if err != nil {
		return status, errors.Wrap(err, "failed to build PersistentVolumeClaim permissions CronJobs")
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPublicFiles, k8ssync.CronJob(drupal, pvcPublicFiles.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPublicFiles.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPublicDirs, k8ssync.CronJob(drupal, pvcPublicDirs.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPublicDirs.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPublicOwner, k8ssync.CronJob(drupal, pvcPublicOwner.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPublicOwner.ObjectMeta.Name, result)

	pvcPrivateFiles, pvcPrivateDirs, pvcPrivateOwner, err := buildVolumePermissions(fmt.Sprintf("%s-private", name), drupal, pvcPrivate, drupal.Spec.Volume.Private)
	if err != nil {
		return status, errors.Wrap(err, "failed to build PersistentVolumeClaim permissions CronJobs")
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPrivateFiles, k8ssync.CronJob(drupal, pvcPrivateFiles.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPrivateFiles.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPrivateDirs, k8ssync.CronJob(drupal, pvcPrivateDirs.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPrivateDirs.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcPrivateOwner, k8ssync.CronJob(drupal, pvcPrivateOwner.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcPrivateOwner.ObjectMeta.Name, result)

	pvcTemporaryFiles, pvcTemporaryDirs, pvcTemporaryOwner, err := buildVolumePermissions(fmt.Sprintf("%s-temporary", name), drupal, pvcTemporary, drupal.Spec.Volume.Temporary)
	if err != nil {
		return status, errors.Wrap(err, "failed to build PersistentVolumeClaim permissions CronJobs")
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcTemporaryFiles, k8ssync.CronJob(drupal, pvcTemporaryFiles.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcTemporaryFiles.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcTemporaryDirs, k8ssync.CronJob(drupal, pvcTemporaryDirs.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcTemporaryDirs.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvcTemporaryOwner, k8ssync.CronJob(drupal, pvcTemporaryOwner.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync CronJob")
	}
	log.Infof("Synced CronJob %s with status: %s", pvcTemporaryOwner.ObjectMeta.Name, result)

	configMapOverride := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-override", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Data: make(map[string]string),
	}

	secretOverride := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-override", name),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		},
		Data: make(map[string][]byte),
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, configMapOverride, k8ssync.ConfigMap(drupal, configMapOverride.Data, false, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync ConfigMap")
	}
	log.Infof("Synced ConfigMap %s with status: %s", configMapOverride.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, secretOverride, k8ssync.Secret(drupal, secretOverride.Data, false, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Secret")
	}
	log.Infof("Synced Secret %s with status: %s", secretOverride.ObjectMeta.Name, result)

	mysqlStatus := make(map[string]appv1beta1.DrupalStatusMySQL)
	mysqlBackup := make(map[string]extensionsv1beta1.BackupSpecMySQL)

	for mysqlKey, mysqlValue := range drupal.Spec.MySQL {
		mysqlMetadata := metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", name, mysqlKey),
			Namespace: drupal.ObjectMeta.Namespace,
			Labels:    commonlabels,
		}

		database := &mysqlv1beta1.Database{
			ObjectMeta: mysqlMetadata,
			Spec: mysqlv1beta1.DatabaseSpec{
				Provisioner: mysqlValue.Class,
				Privileges: []string{
					mysql.PrivilegeAll,
				},
			},
		}

		result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, database, k8ssync.Database(drupal, database.Spec, r.scheme))
		if err != nil {
			return status, errors.Wrap(err, "failed to sync MySQL Database")
		}
		log.Infof("Synced Database %s with status: %s", database.ObjectMeta.Name, result)

		var (
			mysqlKeyHostname = fmt.Sprintf("mysql.%s.hostname", mysqlKey)
			mysqlKeyPort     = fmt.Sprintf("mysql.%s.port", mysqlKey)
			mysqlKeyDatabase = fmt.Sprintf("mysql.%s.database", mysqlKey)
			mysqlKeyUsername = fmt.Sprintf("mysql.%s.username", mysqlKey)
			mysqlKeyPassword = fmt.Sprintf("mysql.%s.password", mysqlKey)
		)

		if database.Status.Connection.Hostname != "" {
			configMapDefault.Data[mysqlKeyHostname] = database.Status.Connection.Hostname
		}
		if database.Status.Connection.Port != 0 {
			configMapDefault.Data[mysqlKeyPort] = strconv.Itoa(database.Status.Connection.Port)
		}
		if database.Status.Connection.Database != "" {
			configMapDefault.Data[mysqlKeyDatabase] = database.Status.Connection.Database
		}

		if database.Status.Connection.Username != "" {
			secretDefault.Data[mysqlKeyUsername] = []byte(database.Status.Connection.Username)
		}
		if database.Status.Connection.Password != "" {
			secretDefault.Data[mysqlKeyPassword] = []byte(database.Status.Connection.Password)
		}

		mysqlStatus[mysqlKey] = appv1beta1.DrupalStatusMySQL{
			ConfigMap: appv1beta1.DrupalStatusMySQLConfigMap{
				Name: configMapDefault.ObjectMeta.Name,
				Keys: appv1beta1.DrupalStatusMySQLConfigMapKeys{
					Hostname: mysqlKeyHostname,
					Port:     mysqlKeyPort,
					Database: mysqlKeyDatabase,
				},
			},
			Secret: appv1beta1.DrupalStatusMySQLSecret{
				Name: secretDefault.ObjectMeta.Name,
				Keys: appv1beta1.DrupalStatusMySQLSecretKeys{
					Username: mysqlKeyUsername,
					Password: mysqlKeyPassword,
				},
			},
		}

		mysqlBackup[mysqlKey] = extensionsv1beta1.BackupSpecMySQL{
			ConfigMap: extensionsv1beta1.BackupSpecMySQLConfigMap{
				Name: configMapDefault.ObjectMeta.Name,
				Keys: extensionsv1beta1.BackupSpecMySQLConfigMapKeys{
					Hostname: mysqlKeyHostname,
					Port:     mysqlKeyPort,
					Database: mysqlKeyDatabase,
				},
			},
			Secret: extensionsv1beta1.BackupSpecMySQLSecret{
				Name: secretDefault.ObjectMeta.Name,
				Keys: extensionsv1beta1.BackupSpecMySQLSecretKeys{
					Username: mysqlKeyUsername,
					Password: mysqlKeyPassword,
				},
			},
		}
	}

	fpmMetadata := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", name, LayerFPM),
		Namespace: drupal.ObjectMeta.Namespace,
		Labels: map[string]string{
			LabelAppName:  drupal.ObjectMeta.Name,
			LabelAppType:  Application,
			LabelAppLayer: LayerFPM,
		},
	}

	var (
		fpmDeploymentGrace = int64(corev1.DefaultTerminationGracePeriodSeconds)
	)

	fpmDeployment := &appsv1.Deployment{
		ObjectMeta: fpmMetadata,
		Spec: appsv1.DeploymentSpec{
			Replicas: &drupal.Spec.FPM.Autoscaling.Replicas.Min,
			Selector: &metav1.LabelSelector{
				MatchLabels: fpmMetadata.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: fpmMetadata.Labels,
					Annotations: map[string]string{
						prometheus.AnnotationScrape: prometheus.ScrapeTrue,
						prometheus.AnnotationPort:   r.exporters.FPM.Port,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            PodContainerFPM,
							Image:           drupal.Spec.FPM.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							SecurityContext: &corev1.SecurityContext{
								ReadOnlyRootFilesystem: &drupal.Spec.FPM.ReadOnly,
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: drupal.Spec.FPM.Port,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							// @todo, LivenessProbe.
							// @todo, Healthz.
							Resources: drupal.Spec.FPM.Resources,
							Env: []corev1.EnvVar{
								generate.EnvVar(PodEnvSkipperEnv, drupal.ObjectMeta.Name),
								generate.EnvVarConfigMap(PodEnvNewRelicApp, drupal.Spec.NewRelic.ConfigMap.Name, configMapOverride, true),
								generate.EnvVarSecret(PodEnvNewRelicLicense, drupal.Spec.NewRelic.Secret.License, secretOverride, true),
								generate.EnvVarConfigMap(PodEnvNewRelicEnabled, drupal.Spec.NewRelic.ConfigMap.Enabled, configMapOverride, true),
							},
							VolumeMounts: []corev1.VolumeMount{
								generate.Mount(VolumeConfigDefault, drupal.Spec.ConfigMap.Default.Path, true),
								generate.Mount(VolumeConfigOverride, drupal.Spec.ConfigMap.Override.Path, true),
								generate.Mount(VolumeSecretDefault, drupal.Spec.Secret.Default.Path, true),
								generate.Mount(VolumeSecretOverride, drupal.Spec.Secret.Override.Path, true),
								generate.Mount(VolumePublic, drupal.Spec.Volume.Public.Path, false),
								generate.Mount(VolumeTemporary, drupal.Spec.Volume.Temporary.Path, false),
								generate.Mount(VolumePrivate, drupal.Spec.Volume.Private.Path, false),
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
						{
							Name:            PodContainerExporter,
							Image:           r.exporters.FPM.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    r.exporters.FPM.CPU,
									corev1.ResourceMemory: r.exporters.FPM.Memory,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    r.exporters.FPM.CPU,
									corev1.ResourceMemory: r.exporters.FPM.Memory,
								},
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					Volumes: []corev1.Volume{
						generate.VolumeConfigMap(VolumeConfigDefault, configMapDefault.ObjectMeta.Name),
						generate.VolumeConfigMap(VolumeConfigOverride, configMapOverride.ObjectMeta.Name),
						generate.VolumeSecret(VolumeSecretDefault, secretDefault.ObjectMeta.Name),
						generate.VolumeSecret(VolumeSecretOverride, secretOverride.ObjectMeta.Name),
						generate.VolumeClaim(VolumePublic, pvcPublic.ObjectMeta.Name),
						generate.VolumeClaim(VolumeTemporary, pvcTemporary.ObjectMeta.Name),
						generate.VolumeClaim(VolumePrivate, pvcPrivate.ObjectMeta.Name),
					},
					// The below are fields which need to be set so we can perform an "deep equal"
					// without always having difference.
					SecurityContext:               &corev1.PodSecurityContext{},
					SchedulerName:                 corev1.DefaultSchedulerName,
					DNSPolicy:                     corev1.DNSClusterFirst,
					TerminationGracePeriodSeconds: &fpmDeploymentGrace,
					RestartPolicy:                 corev1.RestartPolicyAlways,
				},
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, fpmDeployment, k8ssync.Deployment(drupal, fpmDeployment.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Deployment")
	}
	log.Infof("Synced Deployment %s with status: %s", fpmDeployment.ObjectMeta.Name, result)

	fpmHPA := &autoscalingv1.HorizontalPodAutoscaler{
		ObjectMeta: fpmMetadata,
		Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
				APIVersion: appsv1.SchemeGroupVersion.String(),
				Kind:       "Deployment", // @todo, Find a const.
				Name:       fpmDeployment.ObjectMeta.Name,
			},
			MinReplicas:                    &drupal.Spec.FPM.Autoscaling.Replicas.Min,
			MaxReplicas:                    drupal.Spec.FPM.Autoscaling.Replicas.Max,
			TargetCPUUtilizationPercentage: &drupal.Spec.FPM.Autoscaling.Trigger.CPU,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, fpmHPA, k8ssync.HPA(drupal, fpmHPA.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync HorizontalPodAutoscaler")
	}
	log.Infof("Synced HorizontalPodAutoscaler %s with status: %s", fpmHPA.ObjectMeta.Name, result)

	fpmService := &corev1.Service{
		ObjectMeta: fpmMetadata,
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port: drupal.Spec.FPM.Port,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: drupal.Spec.FPM.Port,
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: fpmMetadata.Labels,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, fpmService, k8ssync.Service(drupal, fpmService.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Service")
	}
	log.Infof("Synced Service %s with status: %s", fpmService.ObjectMeta.Name, result)

	nginxMetadata := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", name, LayerNginx),
		Namespace: drupal.ObjectMeta.Namespace,
		Labels: map[string]string{
			LabelAppName:  drupal.ObjectMeta.Name,
			LabelAppType:  Application,
			LabelAppLayer: LayerNginx,
		},
	}

	nginxeploymentGrace := int64(corev1.DefaultTerminationGracePeriodSeconds)

	nginxDeployment := &appsv1.Deployment{
		ObjectMeta: nginxMetadata,
		Spec: appsv1.DeploymentSpec{
			Replicas: &drupal.Spec.Nginx.Autoscaling.Replicas.Min,
			Selector: &metav1.LabelSelector{
				MatchLabels: nginxMetadata.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: nginxMetadata.Labels,
					Annotations: map[string]string{
						prometheus.AnnotationScrape: prometheus.ScrapeTrue,
						prometheus.AnnotationPort:   r.exporters.Nginx.Port,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            PodContainerNginx,
							Image:           drupal.Spec.Nginx.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							SecurityContext: &corev1.SecurityContext{
								ReadOnlyRootFilesystem: &drupal.Spec.Nginx.ReadOnly,
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: drupal.Spec.Nginx.Port,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								generate.Mount(VolumePublic, drupal.Spec.Volume.Public.Path, false),
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
						{
							Name:            PodContainerExporter,
							Image:           r.exporters.Nginx.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    r.exporters.Nginx.CPU,
									corev1.ResourceMemory: r.exporters.Nginx.Memory,
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    r.exporters.Nginx.CPU,
									corev1.ResourceMemory: r.exporters.Nginx.Memory,
								},
							},
							TerminationMessagePath:   corev1.TerminationMessagePathDefault,
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					Volumes: []corev1.Volume{
						generate.VolumeClaim(VolumePublic, pvcPublic.ObjectMeta.Name),
					},
					HostAliases: []corev1.HostAlias{
						{
							IP: fpmService.Spec.ClusterIP,
							Hostnames: []string{
								drupal.Spec.Nginx.HostAlias.FPM,
							},
						},
					},
					// The below are fields which need to be set so we can perform an "deep equal"
					// without always having difference.
					SecurityContext:               &corev1.PodSecurityContext{},
					SchedulerName:                 corev1.DefaultSchedulerName,
					DNSPolicy:                     corev1.DNSClusterFirst,
					TerminationGracePeriodSeconds: &nginxeploymentGrace,
					RestartPolicy:                 corev1.RestartPolicyAlways,
				},
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, nginxDeployment, k8ssync.Deployment(drupal, nginxDeployment.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Deployment")
	}
	log.Infof("Synced Deployment %s with status: %s", nginxDeployment.ObjectMeta.Name, result)

	nginxHPA := &autoscalingv1.HorizontalPodAutoscaler{
		ObjectMeta: nginxMetadata,
		Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
				APIVersion: appsv1.SchemeGroupVersion.String(),
				Kind:       "Deployment", // @todo, Find a const.
				Name:       nginxDeployment.ObjectMeta.Name,
			},
			MinReplicas:                    &drupal.Spec.Nginx.Autoscaling.Replicas.Min,
			MaxReplicas:                    drupal.Spec.Nginx.Autoscaling.Replicas.Max,
			TargetCPUUtilizationPercentage: &drupal.Spec.Nginx.Autoscaling.Trigger.CPU,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, nginxHPA, k8ssync.HPA(drupal, nginxHPA.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync HorizontalPodAutoscaler")
	}
	log.Infof("Synced HorizontalPodAutoscaler %s with status: %s", nginxHPA.ObjectMeta.Name, result)

	nginxService := &corev1.Service{
		ObjectMeta: nginxMetadata,
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port: drupal.Spec.Nginx.Port,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: drupal.Spec.Nginx.Port,
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: nginxMetadata.Labels,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, nginxService, k8ssync.Service(drupal, nginxService.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Service")
	}
	log.Infof("Synced Service %s with status: %s", nginxService.ObjectMeta.Name, result)

	var (
		cronStatus   = make(map[string]appv1beta1.DrupalStatusCron)
		cronDeadLine = int64(1000)
		cronLabels   = map[string]string{
			LabelAppName:  drupal.ObjectMeta.Name,
			LabelAppType:  Application,
			LabelAppLayer: LayerBatch,
		}
	)

	for cronKey, cronValue := range drupal.Spec.Cron {
		cronjob := &batchv1beta1.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", name, cronKey),
				Namespace: drupal.ObjectMeta.Namespace,
				Labels:    cronLabels,
			},
			Spec: batchv1beta1.CronJobSpec{
				Schedule:                   cronValue.Schedule,
				StartingDeadlineSeconds:    &cronDeadLine,
				ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
				SuccessfulJobsHistoryLimit: &cronValue.KeepSuccess,
				FailedJobsHistoryLimit:     &cronValue.KeepFailed,
				JobTemplate: batchv1beta1.JobTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: drupal.ObjectMeta.Namespace,
						Labels:    cronLabels,
					},
					Spec: batchv1.JobSpec{
						BackoffLimit: &cronValue.Retries,
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Namespace: drupal.ObjectMeta.Namespace,
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:            PodContainerFPM,
										Image:           cronValue.Image,
										ImagePullPolicy: corev1.PullIfNotPresent,
										SecurityContext: &corev1.SecurityContext{
											ReadOnlyRootFilesystem: &cronValue.ReadOnly,
										},
										Command: []string{
											"/bin/bash",
											"-c",
										},
										Args: []string{
											cronValue.Command,
										},
										// @todo, ReadOnly
										Resources: cronValue.Resources,
										Env: []corev1.EnvVar{
											generate.EnvVar(PodEnvSkipperEnv, drupal.ObjectMeta.Name),
											generate.EnvVarConfigMap(PodEnvNewRelicApp, drupal.Spec.NewRelic.ConfigMap.Name, configMapOverride, true),
											generate.EnvVarSecret(PodEnvNewRelicLicense, drupal.Spec.NewRelic.Secret.License, secretOverride, true),
											generate.EnvVarConfigMap(PodEnvNewRelicEnabled, drupal.Spec.NewRelic.ConfigMap.Enabled, configMapOverride, true),
										},
										VolumeMounts: []corev1.VolumeMount{
											generate.Mount(VolumeConfigDefault, drupal.Spec.ConfigMap.Default.Path, true),
											generate.Mount(VolumeConfigOverride, drupal.Spec.ConfigMap.Override.Path, true),
											generate.Mount(VolumeSecretDefault, drupal.Spec.Secret.Default.Path, true),
											generate.Mount(VolumeSecretOverride, drupal.Spec.Secret.Override.Path, true),
											generate.Mount(VolumePublic, drupal.Spec.Volume.Public.Path, false),
											generate.Mount(VolumeTemporary, drupal.Spec.Volume.Temporary.Path, false),
											generate.Mount(VolumePrivate, drupal.Spec.Volume.Private.Path, false),
										},
										TerminationMessagePath:   corev1.TerminationMessagePathDefault,
										TerminationMessagePolicy: corev1.TerminationMessageReadFile,
									},
								},
								Volumes: []corev1.Volume{
									generate.VolumeConfigMap(VolumeConfigDefault, configMapDefault.ObjectMeta.Name),
									generate.VolumeConfigMap(VolumeConfigOverride, configMapOverride.ObjectMeta.Name),
									generate.VolumeSecret(VolumeSecretDefault, secretDefault.ObjectMeta.Name),
									generate.VolumeSecret(VolumeSecretOverride, secretOverride.ObjectMeta.Name),
									generate.VolumeClaim(VolumePublic, pvcPublic.ObjectMeta.Name),
									generate.VolumeClaim(VolumeTemporary, pvcTemporary.ObjectMeta.Name),
									generate.VolumeClaim(VolumePrivate, pvcPrivate.ObjectMeta.Name),
								},
								// The below are fields which need to be set so we can perform an "deep equal"
								// without always having difference.
								SecurityContext:               &corev1.PodSecurityContext{},
								SchedulerName:                 corev1.DefaultSchedulerName,
								DNSPolicy:                     corev1.DNSClusterFirst,
								TerminationGracePeriodSeconds: &fpmDeploymentGrace,
								RestartPolicy:                 corev1.RestartPolicyNever,
							},
						},
					},
				},
			},
		}

		result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, cronjob, k8ssync.CronJob(drupal, cronjob.Spec, r.scheme))
		if err != nil {
			return status, errors.Wrap(err, "failed to sync CronJob")
		}
		log.Infof("Synced CronJob %s with status: %s", cronjob.ObjectMeta.Name, result)

		cronStatus[cronKey] = appv1beta1.DrupalStatusCron{
			LastScheduleTime: cronjob.Status.LastScheduleTime,
		}
	}

	// @todo, CronJob cleanup.

	smtp := &extensionsv1beta1.SMTP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: drupal.ObjectMeta.Namespace,
			Labels: map[string]string{
				LabelAppName: drupal.ObjectMeta.Name,
				LabelAppType: Application,
			},
		},
		Spec: extensionsv1beta1.SMTPSpec{
			From: drupal.Spec.SMTP.From,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, smtp, k8ssync.SMTP(drupal, smtp.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync SMTP")
	}
	log.Infof("Synced SMTP %s with status: %s", smtp.ObjectMeta.Name, result)

	if drupal.Spec.SMTP.From.Address != "" {
		configMapDefault.Data["smtp.from.address"] = drupal.Spec.SMTP.From.Address
	}

	if smtp.Status.Connection.Hostname != "" {
		configMapDefault.Data["smtp.hostname"] = smtp.Status.Connection.Hostname
	}

	if smtp.Status.Connection.Port != 0 {
		configMapDefault.Data["smtp.port"] = strconv.Itoa(smtp.Status.Connection.Port)
	}

	if smtp.Status.Connection.Username != "" {
		configMapDefault.Data["smtp.username"] = smtp.Status.Connection.Username
	}

	if smtp.Status.Connection.Password != "" {
		secretDefault.Data["smtp.password"] = []byte(smtp.Status.Connection.Password)
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, configMapDefault, k8ssync.ConfigMap(drupal, configMapDefault.Data, true, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync ConfigMap")
	}
	log.Infof("Synced ConfigMap %s with status: %s", configMapDefault.ObjectMeta.Name, result)

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, secretDefault, k8ssync.Secret(drupal, secretDefault.Data, true, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Secret")
	}
	log.Infof("Synced Secret %s with status: %s", secretDefault.ObjectMeta.Name, result)

	exec := &extensionsv1beta1.Exec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: drupal.ObjectMeta.Namespace,
			Labels: map[string]string{
				LabelAppName: drupal.ObjectMeta.Name,
				LabelAppType: Application,
			},
		},
		Spec: extensionsv1beta1.ExecSpec{
			Entrypoint: PodContainerCLI,
			Template: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            PodContainerCLI,
						Image:           drupal.Spec.Exec.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{
							ReadOnlyRootFilesystem: &drupal.Spec.Exec.ReadOnly,
						},
						// @todo, ReadOnly
						// The sleep command will return an "Exit 0" after the timeout is up.
						// This means the pod will stop and be marked as complete.
						Command: []string{
							"sleep",
							strconv.Itoa(drupal.Spec.Exec.Timeout),
						},
						Resources: drupal.Spec.Exec.Resources,
						Env: []corev1.EnvVar{
							generate.EnvVar(PodEnvSkipperEnv, drupal.ObjectMeta.Name),
							generate.EnvVarConfigMap(PodEnvNewRelicApp, drupal.Spec.NewRelic.ConfigMap.Name, configMapOverride, true),
							generate.EnvVarSecret(PodEnvNewRelicLicense, drupal.Spec.NewRelic.Secret.License, secretOverride, true),
							generate.EnvVarConfigMap(PodEnvNewRelicEnabled, drupal.Spec.NewRelic.ConfigMap.Enabled, configMapOverride, true),
						},
						VolumeMounts: []corev1.VolumeMount{
							generate.Mount(VolumeConfigDefault, drupal.Spec.ConfigMap.Default.Path, true),
							generate.Mount(VolumeConfigOverride, drupal.Spec.ConfigMap.Override.Path, true),
							generate.Mount(VolumeSecretDefault, drupal.Spec.Secret.Default.Path, true),
							generate.Mount(VolumeSecretOverride, drupal.Spec.Secret.Override.Path, true),
							generate.Mount(VolumePublic, drupal.Spec.Volume.Public.Path, false),
							generate.Mount(VolumeTemporary, drupal.Spec.Volume.Temporary.Path, false),
							generate.Mount(VolumePrivate, drupal.Spec.Volume.Private.Path, false),
						},
						TerminationMessagePath:   corev1.TerminationMessagePathDefault,
						TerminationMessagePolicy: corev1.TerminationMessageReadFile,
					},
				},
				Volumes: []corev1.Volume{
					generate.VolumeConfigMap(VolumeConfigDefault, configMapDefault.ObjectMeta.Name),
					generate.VolumeConfigMap(VolumeConfigOverride, configMapOverride.ObjectMeta.Name),
					generate.VolumeSecret(VolumeSecretDefault, secretDefault.ObjectMeta.Name),
					generate.VolumeSecret(VolumeSecretOverride, secretOverride.ObjectMeta.Name),
					generate.VolumeClaim(VolumePublic, pvcPublic.ObjectMeta.Name),
					generate.VolumeClaim(VolumeTemporary, pvcTemporary.ObjectMeta.Name),
					generate.VolumeClaim(VolumePrivate, pvcPrivate.ObjectMeta.Name),
				},
				// The below are fields which need to be set so we can perform an "deep equal"
				// without always having difference.
				SecurityContext:               &corev1.PodSecurityContext{},
				SchedulerName:                 corev1.DefaultSchedulerName,
				DNSPolicy:                     corev1.DNSClusterFirst,
				TerminationGracePeriodSeconds: &fpmDeploymentGrace,
				RestartPolicy:                 corev1.RestartPolicyNever,
			},
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, exec, k8ssync.Exec(drupal, exec.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Exec")
	}
	log.Infof("Synced Exec %s with status: %s", exec.ObjectMeta.Name, result)

	backup := &extensionsv1beta1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: drupal.ObjectMeta.Namespace,
			Labels: map[string]string{
				LabelAppName: drupal.ObjectMeta.Name,
				LabelAppType: Application,
			},
		},
		Spec: extensionsv1beta1.BackupSpec{
			Schedule: drupal.Spec.Backup.Schedule,
			Volumes: map[string]extensionsv1beta1.BackupSpecVolume{
				VolumePublic: {
					ClaimName: pvcPublic.ObjectMeta.Name,
				},
				VolumePrivate: {
					ClaimName: pvcPrivate.ObjectMeta.Name,
				},
			},
			MySQL: mysqlBackup,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, backup, k8ssync.Backup(drupal, backup.Spec, r.scheme))
	if err != nil {
		return status, errors.Wrap(err, "failed to sync Backup")
	}
	log.Infof("Synced Backup %s with status: %s", backup.ObjectMeta.Name, result)

	status = appv1beta1.DrupalStatus{
		ObservedGeneration: drupal.ObjectMeta.Generation,
		Labels: appv1beta1.DrupalStatusLabels{
			All:   commonlabels,
			Nginx: nginxMetadata.Labels,
			FPM:   fpmMetadata.Labels,
			Cron:  cronLabels,
		},
		Nginx: appv1beta1.DrupalStatusNginx{
			Phase:    deploymentutils.GetPhase(nginxDeployment),
			Service:  nginxService.ObjectMeta.Name,
			Image:    drupal.Status.Nginx.Image, // We use the previous image here and update it below if it is "deployed".
			Replicas: nginxDeployment.Status.Replicas,
			// @todo metrics
		},
		FPM: appv1beta1.DrupalStatusFPM{
			Phase:    deploymentutils.GetPhase(fpmDeployment),
			Service:  fpmService.ObjectMeta.Name,
			Image:    drupal.Status.FPM.Image, // We use the previous image here and update it below if it is "deployed".
			Replicas: fpmDeployment.Status.Replicas,
			// @todo metrics
		},
		Volume: appv1beta1.DrupalStatusVolumes{
			Public: appv1beta1.DrupalStatusVolume{
				Name:  pvcPublic.ObjectMeta.Name,
				Phase: pvcPublic.Status.Phase,
			},
			Private: appv1beta1.DrupalStatusVolume{
				Name:  pvcPrivate.ObjectMeta.Name,
				Phase: pvcPrivate.Status.Phase,
			},
			Temporary: appv1beta1.DrupalStatusVolume{
				Name:  pvcTemporary.ObjectMeta.Name,
				Phase: pvcTemporary.Status.Phase,
			},
		},
		MySQL: mysqlStatus,
		Cron:  cronStatus,
		ConfigMap: appv1beta1.DrupalStatusConfigMaps{
			Default: appv1beta1.DrupalStatusConfigMap{
				Name:  configMapDefault.ObjectMeta.Name,
				Count: len(configMapDefault.Data),
			},
			Override: appv1beta1.DrupalStatusConfigMap{
				Name:  configMapOverride.ObjectMeta.Name,
				Count: len(configMapOverride.Data),
			},
		},
		Secret: appv1beta1.DrupalStatusSecrets{
			Default: appv1beta1.DrupalStatusSecret{
				Name:  secretDefault.ObjectMeta.Name,
				Count: len(secretDefault.Data),
			},
			Override: appv1beta1.DrupalStatusSecret{
				Name:  secretOverride.ObjectMeta.Name,
				Count: len(secretOverride.Data),
			},
		},
		Exec: appv1beta1.DrupalStatusExec{
			Name: exec.ObjectMeta.Name,
		},
		SMTP: appv1beta1.DrupalStatusSMTP{
			Verification: smtp.Status.Verification,
		},
		Backup: appv1beta1.DrupalStatusBackup{
			Name:             backup.ObjectMeta.Name,
			LastScheduleTime: backup.Status.LastScheduleTime,
		},
	}

	if status.Nginx.Phase == deploymentutils.PhaseDeployed {
		status.Nginx.Image = drupal.Spec.Nginx.Image
	}

	if status.FPM.Phase == deploymentutils.PhaseDeployed {
		status.FPM.Image = drupal.Spec.FPM.Image
	}

	return status, nil
}

// SyncStatus with the Drupal object.
func (r *ReconcileDrupal) SyncStatus(drupal *appv1beta1.Drupal, status appv1beta1.DrupalStatus) error {
	if reflect.DeepEqual(drupal.Status, status) {
		return nil
	}

	drupal.Status = status

	return r.Client.Status().Update(context.TODO(), drupal)
}
