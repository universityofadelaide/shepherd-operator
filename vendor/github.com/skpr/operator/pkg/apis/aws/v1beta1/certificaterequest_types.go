package v1beta1

import (
	"crypto/md5"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertificateRequestSpec defines the desired state of CertificateRequest
type CertificateRequestSpec struct {
	// Primary domain for the certificate request.
	CommonName string `json:"commonName"`
	// Additional domains for the certificate request.
	AlternateNames []string `json:"alternateNames,omitempty"`
}

// Hash value derived from the CommonName and AlternativeNames.
func (s CertificateRequestSpec) Hash() string {
	domains := append(s.AlternateNames, s.CommonName)
	val := strings.Join(domains, "")
	return fmt.Sprintf("%x", md5.Sum([]byte(val)))
}

// CertificateRequestStatus defines the observed state of CertificateRequest
type CertificateRequestStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Machine identifier for the certificate request.
	ARN string `json:"arn,omitempty"`
	// Domain list for the certificate.
	Domains []string `json:"domains,omitempty"`
	// Current state of the certificate eg. ISSUED.
	State string `json:"state,omitempty"`
	// Details used to validate a certificate request.
	Validate []ValidateRecord `json:"validate,omitempty"`
}

// ValidateRecord provide details to site administrators on how to validate a certificate.
type ValidateRecord struct {
	// The name of DNS validation record.
	Name string `json:"name,omitempty"`
	// The type of DNS validation record.
	Type string `json:"type,omitempty"`
	// The value of DNS validation record.
	Value string `json:"value,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertificateRequest is the Schema for the certificaterequests API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CertificateRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateRequestSpec   `json:"spec,omitempty"`
	Status CertificateRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertificateRequestList contains a list of CertificateRequest
type CertificateRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertificateRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CertificateRequest{}, &CertificateRequestList{})
}
