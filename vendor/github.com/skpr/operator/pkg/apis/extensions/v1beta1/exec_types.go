package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExecSpec defines the desired state of Exec
type ExecSpec struct {
	// Container which commands will be executed.
	Entrypoint string `json:"entrypoint"`
	// Template used when provisioning an execution environment.
	Template corev1.PodSpec `json:"template"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Exec is the Schema for the execs API
// +k8s:openapi-gen=true
type Exec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ExecSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExecList contains a list of Exec
type ExecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Exec `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Exec{}, &ExecList{})
}
