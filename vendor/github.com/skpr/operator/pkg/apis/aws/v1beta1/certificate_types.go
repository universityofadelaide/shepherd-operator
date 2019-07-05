package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertificateSpec defines the desired state of Certificate
type CertificateSpec struct {
	// Information which will be used to provision a certificate.
	Request CertificateRequestSpec `json:"request"`
}

// CertificateStatus defines the observed state of Certificate
type CertificateStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// The status of the desired certificate.
	Desired CertificateRequestReference `json:"desired,omitempty"`
	// The status of the most recently ISSUED certificate.
	Active CertificateRequestReference `json:"active,omitempty"`
	// Status of all the certificate requests.
	Requests []CertificateRequestReference `json:"requests,omitempty"`
}

// CertificateRequestReference defines the observed state of Certificate
type CertificateRequestReference struct {
	// Reference name for the certificate request.
	Name string `json:"name,omitempty"`
	// Details of the certificate.
	Details CertificateRequestStatus `json:"details,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Certificate is the Schema for the certificates API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Certificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateSpec   `json:"spec,omitempty"`
	Status CertificateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertificateList contains a list of Certificate
type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Certificate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Certificate{}, &CertificateList{})
}
