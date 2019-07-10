package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metav1_shepherd "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

// BackupScheduledSpec defines the desired state of BackupScheduled
type BackupScheduledSpec struct {
	// Schedule is the crontab statement which defines how often a backup should run.
	Schedule string `json:"schedule"`
	// Volumes which will be backed up.
	Volumes map[string]metav1_shepherd.SpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]metav1_shepherd.SpecMySQL `json:"mysql,omitempty"`
}

// BackupScheduledStatus defines the observed state of BackupScheduled
type BackupScheduledStatus struct {
	LastExecutedTime *metav1.Time          `json:"lastExecutedTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupScheduled is the Schema for the backupscheduleds API
// +k8s:openapi-gen=true
type BackupScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduledSpec   `json:"spec,omitempty"`
	Status BackupScheduledStatus `json:"status,omitempty"`
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
