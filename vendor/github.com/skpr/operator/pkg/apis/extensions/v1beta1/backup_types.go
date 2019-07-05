package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// Frequency which the application will be backed up.
	Schedule string `json:"schedule"`
	// Volumes which will be backed up.
	Volumes map[string]BackupSpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]BackupSpecMySQL `json:"mysql,omitempty"`
}

// BackupSpecVolume defines how to backup volumes.
type BackupSpecVolume struct {
	// ClaimName which will be backed up.
	ClaimName string `json:"claimName"`
}

// BackupSpecMySQL defines how to backup MySQL.
type BackupSpecMySQL struct {
	// ConfigMap which will be used for connectivity.
	ConfigMap BackupSpecMySQLConfigMap `json:"configmap"`
	// Secret which will be used for connectivity.
	Secret BackupSpecMySQLSecret `json:"secret"`
}

// BackupSpecMySQLConfigMap defines connection information for MySQL.
type BackupSpecMySQLConfigMap struct {
	// Name of the ConfigMap.
	Name string `json:"name"`
	// ConfigMap keys used for connectivity.
	Keys BackupSpecMySQLConfigMapKeys `json:"keys"`
}

// BackupSpecMySQLConfigMapKeys defines ConfigMap keys for MySQL connectivity.
type BackupSpecMySQLConfigMapKeys struct {
	// Key which was applied to the application for database connectivity.
	Database string `json:"database"`
	// Key which was applied to the application for database connectivity.
	Hostname string `json:"hostname"`
	// Key which was applied to the application for database connectivity.
	Port string `json:"port"`
}

// BackupSpecMySQLSecret defines connection information for MySQL.
type BackupSpecMySQLSecret struct {
	// Name of the Secret.
	Name string `json:"name"`
	// ConfigMap keys used for connectivity.
	Keys BackupSpecMySQLSecretKeys `json:"keys"`
}

// BackupSpecMySQLSecretKeys defines Secret keys for MySQL connectivity.
type BackupSpecMySQLSecretKeys struct {
	// Key which was applied to the application for database connectivity.
	Username string `json:"username"`
	// Key which was applied to the application for database connectivity.
	Password string `json:"password"`
}

// BackupStatus defines the observed state of Backup
type BackupStatus struct {
	// Last time a backup was executed.
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Backup is the Schema for the backups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupList contains a list of Backup
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
