package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

// BackupScheduledSpec defines the desired state of BackupScheduled
type BackupScheduledSpec struct {
	// Retention defines how backup retention behavior.
	Retention shpmetav1.RetentionSpec `json:"retention"`
	// Schedule is the crontab statement which defines how often a backup should run.
	Schedule shpmetav1.ScheduledSpec `json:"schedule"`
	// Volumes which will be backed up.
	Volumes map[string]shpmetav1.SpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]shpmetav1.SpecMySQL `json:"mysql,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupScheduled is the Schema for the backupscheduleds API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type BackupScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduledSpec       `json:"spec,omitempty"`
	Status shpmetav1.ScheduledStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupScheduledList contains a list of BackupScheduled
type BackupScheduledList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupScheduled `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupScheduled{}, &BackupScheduledList{})
}
