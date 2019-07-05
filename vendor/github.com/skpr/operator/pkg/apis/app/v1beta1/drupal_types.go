package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	deploymentutils "github.com/skpr/operator/pkg/utils/k8s/deployment"
)

// DrupalSpec defines the desired state of Drupal
type DrupalSpec struct {
	// Configuration for the Nginx deployment eg. image / resources / scaling.
	Nginx DrupalSpecNginx `json:"nginx"`
	// Configuration for the FPM deployment eg. image / resources / scaling.
	FPM DrupalSpecFPM `json:"fpm"`
	// Configuration for the Execution environment eg. image / resources / timeout.
	Exec DrupalSpecExec `json:"exec"`
	// Volumes which are provisioned for the Drupal application.
	Volume DrupalSpecVolumes `json:"volume"`
	// Database provisioned as part of the application eg. "default" and "migrate".
	MySQL map[string]DrupalSpecMySQL `json:"mysql"`
	// Background tasks which are executed periodically eg. "drush cron"
	Cron map[string]DrupalSpecCron `json:"cron"`
	// Configuration which is exposed to the Drupal application eg. database hostname.
	ConfigMap DrupalSpecConfigMaps `json:"configmap"`
	// Secrets which are exposed to the Drupal application eg. database credentials.
	Secret DrupalSpecSecrets `json:"secret"`
	// NewRelic configuration for performance and debugging.
	NewRelic DrupalSpecNewRelic `json:"newrelic"`
	// SMTP configuration for outbound email.
	SMTP DrupalSpecSMTP `json:"smtp"`
	// Backup configuration for recovery.
	Backup DrupalSpecBackup `json:"backup"`
	// Prometheus configuration for https://www.drupal.org/project/prometheus_exporter.
	Prometheus DrupalSpecPrometheus `json:"prometheus"`
}

// DrupalSpecNginx provides a specification for the Nginx layer.
type DrupalSpecNginx struct {
	// Image which will be rolled out for the deployment.
	Image string `json:"image"`
	// Port which Nginx is running on.
	Port int32 `json:"port"`
	// Marks the filesystem as readonly to ensure code cannot be changed.
	ReadOnly bool `json:"readOnly"`
	// Resource constraints which will be applied to the deployment.
	Resources corev1.ResourceRequirements `json:"resources"`
	// Autoscaling rules which will be applied to the deployment.
	Autoscaling DrupalSpecNginxAutoscaling `json:"autoscaling"`
	// HostAlias which connects Nginx to the FPM service.
	HostAlias DrupalSpecNginxHostAlias `json:"hostAlias"`
}

// DrupalSpecNginxAutoscaling contain autoscaling rules for the Nginx layer.
type DrupalSpecNginxAutoscaling struct {
	// Threshold which will trigger autoscaling events.
	Trigger DrupalSpecNginxAutoscalingTrigger `json:"trigger"`
	// How many copies of the deployment which should be running.
	Replicas DrupalSpecNginxAutoscalingReplicas `json:"replicas"`
}

// DrupalSpecNginxAutoscalingTrigger contain autoscaling triggers for the Nginx layer.
type DrupalSpecNginxAutoscalingTrigger struct {
	// CPU threshold which will trigger autoscaling events.
	CPU int32 `json:"cpu"`
}

// DrupalSpecNginxAutoscalingReplicas contain autoscaling replica rules for the Nginx layer.
type DrupalSpecNginxAutoscalingReplicas struct {
	// Minimum number of replicas for the deployment.
	Min int32 `json:"min"`
	// Maximum number of replicas for the deployment.
	Max int32 `json:"max"`
}

// DrupalSpecNginxHostAlias declares static DNS entries which wire up the Nginx layer to the FPM lay
type DrupalSpecNginxHostAlias struct {
	// HostAlias for the PHP FPM service.
	FPM string `json:"fpm"`
}

// DrupalSpecFPM provides a specification for the FPM layer.
type DrupalSpecFPM struct {
	// Image which will be rolled out for the deployment.
	Image string `json:"image"`
	// Port which PHP-FPM is running on.
	Port int32 `json:"port"`
	// Marks the filesystem as readonly to ensure code cannot be changed.
	ReadOnly bool `json:"readOnly"`
	// Resource constraints which will be applied to the deployment.
	Resources corev1.ResourceRequirements `json:"resources"`
	// Autoscaling rules which will be applied to the deployment.
	Autoscaling DrupalSpecFPMAutoscaling `json:"autoscaling"`
}

// DrupalSpecFPMAutoscaling contain autoscaling rules for the FPM layer.
type DrupalSpecFPMAutoscaling struct {
	// Threshold which will trigger autoscaling events.
	Trigger DrupalSpecFPMAutoscalingTrigger `json:"trigger"`
	// How many copies of the deployment which should be running.
	Replicas DrupalSpecFPMAutoscalingReplicas `json:"replicas"`
}

// DrupalSpecFPMAutoscalingTrigger contain autoscaling triggers for the FPM layer.
type DrupalSpecFPMAutoscalingTrigger struct {
	// CPU threshold which will trigger autoscaling events.
	CPU int32 `json:"cpu"`
}

// DrupalSpecFPMAutoscalingReplicas contain autoscaling replica rules for the FPM layer.
type DrupalSpecFPMAutoscalingReplicas struct {
	// Minimum number of replicas for the deployment.
	Min int32 `json:"min"`
	// Maximum number of replicas for the deployment.
	Max int32 `json:"max"`
}

// DrupalSpecVolumes which will be mounted for this Drupal.
type DrupalSpecVolumes struct {
	// Public filesystem which is mounted into the Nginx and PHP-FPM deployments.
	Public DrupalSpecVolume `json:"public"`
	// Private filesystem which is only mounted into the PHP-FPM deployment.
	Private DrupalSpecVolume `json:"private"`
	// Temporary filesystem which is only mounted into the PHP-FPM deployment.
	Temporary DrupalSpecVolume `json:"temporary"`
}

// DrupalSpecVolume which will be mounted for this Drupal.
type DrupalSpecVolume struct {
	// Path which the volume will be mounted.
	Path string `json:"path"`
	// StorageClass which will be used to provision storage.
	Class string `json:"class"`
	// Amount of storage which will be provisioned.
	Amount string `json:"amount"`
	// Permmissions which will be enforced for this volume.
	Permissions DrupalSpecVolumePermissions `json:"permissions"`
}

// DrupalSpecVolumePermissions describes permisions which will be enforced for a volume.
type DrupalSpecVolumePermissions struct {
	// User name which will be enforced for the volume. Can also be a UID.
	User string `json:"user"`
	// Group name which will be enforced for the volume. Can also be a GID.
	Group string `json:"group"`
	// Directory chmod which will be applied to the volume.
	Directory int `json:"directory"`
	// File chmod which will be applied to the volume.
	File int `json:"file"`
	// CronJob which will be run to enforce permissions.
	CronJob DrupalSpecVolumePermissionsCronJob `json:"cronJob"`
}

// DrupalSpecVolumePermissionsCronJob describes how the permissions will be enforced.
type DrupalSpecVolumePermissionsCronJob struct {
	// Privileged image which will be used when running chown and chmod commands. Needs be an image which runs as root to enforce permissions.
	Image string `json:"image"`
	// Schedule which permissions will be checked.
	Schedule string `json:"schedule"`
	// How long before a Job should be marked as failed if it does not get scheduled in time.
	Deadline int64 `json:"deadline,omitempty"`
	// How many times to get a "Successful" execution before failing.
	Retries int32 `json:"retries,omitempty"`
	// How many past successes to keep.
	KeepSuccess int32 `json:"keepSuccess,omitempty"`
	// How many past failures to keep.
	KeepFailed int32 `json:"keepFailed,omitempty"`
	// Make the root filesystem of this image ready only.
	ReadOnly bool `json:"readOnly,omitempty"`
}

// DrupalSpecMySQL configuration for this Drupal.
type DrupalSpecMySQL struct {
	// Database class which will be used when provisioning the database.
	Class string `json:"class"`
}

// DrupalSpecCron configures a background task for this Drupal.
type DrupalSpecCron struct {
	// Image which will be executed for a background task.
	Image string `json:"image"`
	// Marks the filesystem as readonly to ensure code cannot be changed.
	ReadOnly bool `json:"readOnly"`
	// Command which will be executed for the background task.
	Command string `json:"command"`
	// Schedule whic the background task will be executed.
	Schedule string `json:"schedule"`
	// Resource constraints which will be applied to the background task.
	Resources corev1.ResourceRequirements `json:"resources"`
	// How many times to execute before marking the task as a failure.
	Retries int32 `json:"retries"`
	// How many successful builds to keep.
	KeepSuccess int32 `json:"keepSuccess"`
	// How many failed builds to keep.
	KeepFailed int32 `json:"keepFailed"`
}

// DrupalSpecConfigMaps defines the spec for all configmaps.
type DrupalSpecConfigMaps struct {
	// Configuration which is automatically set.
	Default DrupalSpecConfigMap `json:"default"`
	// Configuration which is user provided.
	Override DrupalSpecConfigMap `json:"override"`
}

// DrupalSpecConfigMap defines the spec for a specific configmap.
type DrupalSpecConfigMap struct {
	// Path which the configuration will be mounted.
	Path string `json:"path"`
}

// DrupalSpecSecrets defines the spec for all secrets.
type DrupalSpecSecrets struct {
	// Secrets which are automatically set.
	Default DrupalSpecSecret `json:"default"`
	// Secrets which are user provided.
	Override DrupalSpecSecret `json:"override"`
}

// DrupalSpecSecret defines the spec for a specific secret.
type DrupalSpecSecret struct {
	// Path which the secrets will be mounted.
	Path string `json:"path"`
}

// DrupalSpecExec provides configuration for an exec template.
type DrupalSpecExec struct {
	// Image which will be used for developer command line access.
	Image string `json:"image"`
	// Marks the filesystem as readonly to ensure code cannot be changed.
	ReadOnly bool `json:"readOnly"`
	// Resource constraints which will be applied to the background task.
	Resources corev1.ResourceRequirements `json:"resources"`
	// How long for the command line environment to exist.
	Timeout int `json:"timeout"`
}

// DrupalSpecNewRelic is used for configuring New Relic configuration.
type DrupalSpecNewRelic struct {
	// ConfigMap which contains New Relic configuration.
	ConfigMap DrupalSpecNewRelicConfigMap `json:"configmap"`
	// Secret which contains New Relic configuration.
	Secret DrupalSpecNewRelicSecret `json:"secret"`
}

// DrupalSpecNewRelicConfigMap for loading New Relic config from a ConfigMap.
type DrupalSpecNewRelicConfigMap struct {
	// Key which determines if New Relic is enabled.
	Enabled string `json:"enabled"`
	// Name of the New Relic application.
	Name string `json:"name"`
}

// DrupalSpecNewRelicSecret for loading New Relic config from a Secret.
type DrupalSpecNewRelicSecret struct {
	// License (API Key) for interacting with New Relic.
	License string `json:"license"`
}

// DrupalSpecSMTP configuration for outbound email.
type DrupalSpecSMTP struct {
	// Configuration for validationing FROM addresses.
	From extensionsv1beta1.SMTPSpecFrom `json:"from"`
}

// DrupalSpecBackup configuration for recovery.
type DrupalSpecBackup struct {
	Schedule string `json:"schedule"`
}

// DrupalSpecPrometheus configures application metrics.
type DrupalSpecPrometheus struct {
	// Token which Prometheus uses to access https://www.drupal.org/project/prometheus_exporter
	Token string `json:"token"`
}

// DrupalStatus defines the observed state of Drupal
type DrupalStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Labels for querying Drupal services.
	Labels DrupalStatusLabels `json:"labels,omitempty"`
	// Status for the Nginx deployment.
	Nginx DrupalStatusNginx `json:"nginx,omitempty"`
	// Status for the PHP-FPM deployment.
	FPM DrupalStatusFPM `json:"fpm,omitempty"`
	// Volume status information.
	Volume DrupalStatusVolumes `json:"volume,omitempty"`
	// MySQL status information.
	MySQL map[string]DrupalStatusMySQL `json:"mysql,omitempty"`
	// Background task information eg. Last execution.
	Cron map[string]DrupalStatusCron `json:"cron,omitempty"`
	// Configuration status.
	ConfigMap DrupalStatusConfigMaps `json:"configmap,omitempty"`
	// Secrets status.
	Secret DrupalStatusSecrets `json:"secret,omitempty"`
	// Execution environment status.
	Exec DrupalStatusExec `json:"exec,omitempty"`
	// SMTP service status.
	SMTP DrupalStatusSMTP `json:"smtp,omitempty"`
	// Backup status.
	Backup DrupalStatusBackup `json:"backup,omitempty"`
}

// DrupalStatusLabels which are used for querying application components.
type DrupalStatusLabels struct {
	// Used to query all Drupal application pods.
	All map[string]string `json:"all,omitempty"`
	// Used to query all Nginx pods.
	Nginx map[string]string `json:"nginx,omitempty"`
	// Used to query all PHP-FPM pods.
	FPM map[string]string `json:"fpm,omitempty"`
	// Used to query all background task pods.
	Cron map[string]string `json:"cron,omitempty"`
}

// DrupalStatusNginx identifies all deployment related status.
type DrupalStatusNginx struct {
	// Current phase of the Nginx deployment eg. InProgress.
	Phase deploymentutils.Phase `json:"phase,omitempty"`
	// Service for routing traffic.
	Service string `json:"service,omitempty"`
	// Current image which has been rolled out.
	Image string `json:"image,omitempty"`
	// Current number of replicas.
	Replicas int32 `json:"replicas,omitempty"`
	// Application metrics.
	Metrics DrupalStatusNginxMetrics `json:"metrics,omitempty"`
}

// DrupalStatusNginxMetrics identifies all nginx metric related status.
type DrupalStatusNginxMetrics struct {
	// Current CPU metric for the Nginx deployment.
	CPU int32 `json:"cpu,omitempty"`
}

// DrupalStatusFPM identifies all deployment related status.
type DrupalStatusFPM struct {
	// Current phase of the PHP-FPM deployment eg. InProgress.
	Phase deploymentutils.Phase `json:"phase,omitempty"`
	// Service for routing traffic.
	Service string `json:"service,omitempty"`
	// Current image which has been rolled out.
	Image string `json:"image,omitempty"`
	// Current number of replicas.
	Replicas int32 `json:"replicas,omitempty"`
	// Application metrics.
	Metrics DrupalStatusFPMMetrics `json:"metrics,omitempty"`
}

// DrupalStatusFPMMetrics identifies all fpm metric related status.
type DrupalStatusFPMMetrics struct {
	// Current CPU metric for the PHP-FPM deployment.
	CPU int32 `json:"cpu,omitempty"`
}

// DrupalStatusVolumes identifies all volume related status.
type DrupalStatusVolumes struct {
	// Current state of the public filesystem volume.
	Public DrupalStatusVolume `json:"public,omitempty"`
	// Current state of the private filesystem volume.
	Private DrupalStatusVolume `json:"private,omitempty"`
	// Current state of the temporary filesystem volume.
	Temporary DrupalStatusVolume `json:"temporary,omitempty"`
}

// DrupalStatusVolume identifies specific volume related status.
type DrupalStatusVolume struct {
	// Name of the volume.
	Name string `json:"name,omitempty"`
	// Current state of the volume.
	Phase corev1.PersistentVolumeClaimPhase `json:"phase,omitempty"`
}

// DrupalStatusMySQL identifies all mysql related status.
type DrupalStatusMySQL struct {
	// Status of the application configuration.
	ConfigMap DrupalStatusMySQLConfigMap `json:"configmap,omitempty"`
	// Status of the application secrets.
	Secret DrupalStatusMySQLSecret `json:"secret,omitempty"`
}

// DrupalStatusMySQLConfigMap identifies all mysql configmap related status.
type DrupalStatusMySQLConfigMap struct {
	// Name of the configmap.
	Name string `json:"name,omitempty"`
	// Keys which can be used for discovery.
	Keys DrupalStatusMySQLConfigMapKeys `json:"keys,omitempty"`
}

// DrupalStatusMySQLConfigMapKeys identifies all mysql configmap keys.
type DrupalStatusMySQLConfigMapKeys struct {
	// Key which was applied to the application for database connectivity.
	Database string `json:"database,omitempty"`
	// Key which was applied to the application for database connectivity.
	Hostname string `json:"hostname,omitempty"`
	// Key which was applied to the application for database connectivity.
	Port string `json:"port,omitempty"`
}

// DrupalStatusMySQLSecret identifies all mysql secret related status.
type DrupalStatusMySQLSecret struct {
	// Name of the secret.
	Name string `json:"name,omitempty"`
	// Keys which can be used for discovery.
	Keys DrupalStatusMySQLSecretKeys `json:"keys,omitempty"`
}

// DrupalStatusMySQLSecretKeys identifies all mysql secret keys.
type DrupalStatusMySQLSecretKeys struct {
	// Key which was applied to the application for database connectivity.
	Username string `json:"username,omitempty"`
	// Key which was applied to the application for database connectivity.
	Password string `json:"password,omitempty"`
}

//DrupalStatusCron identifies all cron related status.
type DrupalStatusCron struct {
	// Last time a background task was executed.
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

// DrupalStatusConfigMaps identifies all config map related status.
type DrupalStatusConfigMaps struct {
	// Status of the generated configuration.
	Default DrupalStatusConfigMap `json:"default,omitempty"`
	// Status of the user provided configuration.
	Override DrupalStatusConfigMap `json:"override,omitempty"`
}

// DrupalStatusConfigMap identifies specific config map related status.
type DrupalStatusConfigMap struct {
	// Name of the configmap.
	Name string `json:"name,omitempty"`
	// How many configuration are assigned to the application.
	Count int `json:"count,omitempty"`
}

// DrupalStatusSecrets identifies all secret related status.
type DrupalStatusSecrets struct {
	// Status of the generated secrets.
	Default DrupalStatusSecret `json:"default,omitempty"`
	// Status of the user provided secrets.
	Override DrupalStatusSecret `json:"override,omitempty"`
}

// DrupalStatusSecret identifies specific secret related status.
type DrupalStatusSecret struct {
	// Name of the secret.
	Name string `json:"name,omitempty"`
	// How many secrets are assigned to the application.
	Count int `json:"count,omitempty"`
}

// DrupalStatusExec identifies specific template which can be loaded by external sources.
type DrupalStatusExec struct {
	// Name of the command line environment template.
	Name string `json:"name,omitempty"`
}

// DrupalStatusSMTP provides the status for the SMTP service.
type DrupalStatusSMTP struct {
	Verification extensionsv1beta1.SMTPStatusVerification `json:"verification,omitempty"`
}

// DrupalStatusBackup provides the status for the Backup.
type DrupalStatusBackup struct {
	Name             string       `json:"name,omitempty"`
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Drupal is the Schema for the drupals API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Drupal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DrupalSpec   `json:"spec,omitempty"`
	Status DrupalStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DrupalList contains a list of Drupal
type DrupalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Drupal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Drupal{}, &DrupalList{})
}
