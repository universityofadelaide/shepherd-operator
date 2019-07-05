package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CloudFrontInvalidationSpec defines the desired state of CloudFrontInvalidation
type CloudFrontInvalidationSpec struct {
	// Name of the CloudFront object.
	Distribution string `json:"distribution"`
	// Paths which to invalidate.
	Paths []string `json:"paths,omitempty"`
}

const (
	// CloudFrontInvalidationCompleted for identifying with a invalidation has completed.
	CloudFrontInvalidationCompleted = "Completed"
)

// CloudFrontInvalidationStatus defines the observed state of CloudFrontInvalidation
type CloudFrontInvalidationStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Machine identifier for querying an invalidation request.
	ID string `json:"id,omitempty"`
	// When the invalidation request was lodged.
	Created string `json:"created,omitempty"`
	// Current state of the invalidation request.
	State string `json:"state,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFrontInvalidation is the Schema for the cloudfrontinvalidations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CloudFrontInvalidation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudFrontInvalidationSpec   `json:"spec,omitempty"`
	Status CloudFrontInvalidationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFrontInvalidationList contains a list of CloudFrontInvalidation
// +kubebuilder:subresource:status
type CloudFrontInvalidationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudFrontInvalidation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudFrontInvalidation{}, &CloudFrontInvalidationList{})
}
