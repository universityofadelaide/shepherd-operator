package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SMTPSpec defines the desired state of SMTP
type SMTPSpec struct {
	// From defines what an application is allowed to send from.
	From SMTPSpecFrom `json:"from"`
}

// SMTPSpecFrom defines what an application is allowed to send from.
type SMTPSpecFrom struct {
	// Address which an application is allowed to send from.
	Address string `json:"address"`
}

// SMTPStatus defines the observed state of SMTP
type SMTPStatus struct {
	// Provides the status of verifying FROM attributes.
	Verification SMTPStatusVerification `json:"verification,omitempty"`
	// Provides connection details for sending email.
	Connection SMTPStatusConnection `json:"connection,omitempty"`
}

// SMTPStatusVerification provides the status of verifying FROM attributes.
type SMTPStatusVerification struct {
	// Address which an application is allowed to send from.
	Address string `json:"address,omitempty"`
}

// SMTPStatusConnection provides connection details for sending email.
type SMTPStatusConnection struct {
	// Hostname used when connecting to the SMTP server.
	Hostname string `json:"hostname,omitempty"`
	// Port used when connecting to the SMTP server.
	Port int `json:"port,omitempty"`
	// Username used when connecting to the SMTP server.
	Username string `json:"username,omitempty"`
	// Password used when connecting to the SMTP server.
	Password string `json:"password,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SMTP is the Schema for the smtps API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type SMTP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SMTPSpec   `json:"spec,omitempty"`
	Status SMTPStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SMTPList contains a list of SMTP
type SMTPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SMTP `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SMTP{}, &SMTPList{})
}
