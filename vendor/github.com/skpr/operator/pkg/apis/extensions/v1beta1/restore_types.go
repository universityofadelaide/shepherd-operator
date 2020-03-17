package v1beta1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
)

// RestoreSpec defines the desired state of Restore
type RestoreSpec struct {
	// BackupName to restore from.
	BackupName string `json:"backupName"`
	// Secret which defines information about the restic secret.
	Secret BackupSpecSecret `json:"secret"`
	// Volumes which will be restored to.
	Volumes map[string]BackupSpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be restored to.
	MySQL map[string]BackupSpecMySQL `json:"mysql,omitempty"`
}

// RestoreStatus defines the observed state of Restore
type RestoreStatus struct {
	// The time the restore started.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// The time the restore completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
	// The phase the restore is in.
	Phase skprmetav1.Phase `json:"phase"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Restore is the Schema for the restores API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Restore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestoreSpec   `json:"spec,omitempty"`
	Status RestoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RestoreList contains a list of Restore
type RestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Restore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Restore{}, &RestoreList{})
}

// GetDuration returns the duration of the backup.
func (r Restore) GetDuration() *time.Duration {
	if r.Status.StartTime != nil && r.Status.CompletionTime != nil {
		duration := r.Status.CompletionTime.Sub(r.Status.StartTime.Time)
		return &duration
	}

	return nil
}
