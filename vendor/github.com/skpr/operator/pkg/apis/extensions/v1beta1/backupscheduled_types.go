package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
)

// BackupScheduledSpec defines the desired state of BackupScheduled
type BackupScheduledSpec struct {
	Schedule skprmetav1.ScheduledSpec `json:"schedule"`
	// Template which will be used to trigger a backup build.
	Template BackupSpec `json:"template"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupScheduled is the Schema for the backupscheduleds API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type BackupScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduledSpec        `json:"spec,omitempty"`
	Status skprmetav1.ScheduledStatus `json:"status,omitempty"`
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
