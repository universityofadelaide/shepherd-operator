package v1beta1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	"github.com/skpr/operator/pkg/utils/slice"
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// Secret which defines information about the restic secret.
	Secret BackupSpecSecret `json:"secret"`
	// Volumes which will be backed up.
	Volumes map[string]BackupSpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]BackupSpecMySQL `json:"mysql,omitempty"`
	// Tags to apply to the restic backup.
	Tags []string `json:"tags,omitempty"`
}

// BackupSpecSecret defines information for the restic secret.
type BackupSpecSecret struct {
	// Name of the secret.
	Name string `json:"name"`
	// Key of the secret data to use for the password
	Key string `json:"key"`
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
	// Database which was applied to the application for database connectivity.
	Database string `json:"database"`
	// Hostname which was applied to the application for database connectivity.
	Hostname string `json:"hostname"`
	// Port which was applied to the application for database connectivity.
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
	// The time the backup started.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// The time the backup completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// The restic id for the backup.
	ResticID string `json:"resticId,omitempty"`
	// The phase the backup is in.
	Phase skprmetav1.Phase `json:"phase"`
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

// GetDuration returns the duration of the backup.
func (b Backup) GetDuration() *time.Duration {
	if b.Status.StartTime != nil && b.Status.CompletionTime != nil {
		duration := b.Status.CompletionTime.Sub(b.Status.StartTime.Time)
		return &duration
	}

	return nil
}

// HasTag checks if this backup has a tag.
func (b Backup) HasTag(tag string) bool {
	return slice.Contains(b.Spec.Tags, tag)
}
