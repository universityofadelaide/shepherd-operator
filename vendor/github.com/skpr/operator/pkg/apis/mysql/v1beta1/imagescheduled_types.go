package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
)

// ImageScheduledSpec defines the desired state of ImageScheduled
type ImageScheduledSpec struct {
	Schedule skprmetav1.ScheduledSpec `json:"schedule"`
	// Template which will be used to trigger an image build.
	Template ImageSpec `json:"template"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageScheduled is the Schema for the imagescheduleds API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ImageScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageScheduledSpec         `json:"spec,omitempty"`
	Status skprmetav1.ScheduledStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageScheduledList contains a list of ImageScheduled
type ImageScheduledList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageScheduled `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageScheduled{}, &ImageScheduledList{})
}
