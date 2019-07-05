package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CloudFrontStateInProgress is used for determining if a CloudFront is rolling out.
	CloudFrontStateInProgress = "InProgress"
	// CloudFrontStateDeployed is used for determining if a CloudFront has finished deploying.
	CloudFrontStateDeployed = "Deployed"
)

// CloudFrontSpec defines the desired state of CloudFront
type CloudFrontSpec struct {
	// Aliases which CloudFront will respond to.
	Aliases []string `json:"aliases"`
	// Certificate which is applied to this CloudFront distribution.
	Certificate CloudFrontSpecCertificate `json:"certificate,omitempty"`
	// Firewall configuration for this CloudFront distribution.
	Firewall CloudFrontSpecFirewall `json:"firewall,omitempty"`
	// Behavior applied to this CloudFront distribution eg. Headers and Cookies.
	Behavior CloudFrontSpecBehavior `json:"behavior"`
	// Information CloudFront uses to connect to the backend.
	Origin CloudFrontSpecOrigin `json:"origin"`
}

// CloudFrontSpecCertificate declares a certificate to use for encryption.
type CloudFrontSpecCertificate struct {
	// Machine identifier for referencing a certificate.
	ARN string `json:"arn"`
}

// CloudFrontSpecFirewall declares a firewall which this CloudFront is associated with.
type CloudFrontSpecFirewall struct {
	// Machine identifier for referencing a firewall.
	ARN string `json:"arn"`
}

// CloudFrontSpecBehavior declares the behaviour which will be applied to this CloudFront distribution.
type CloudFrontSpecBehavior struct {
	// Whitelist of headers and cookies.
	Whitelist CloudFrontSpecBehaviorWhitelist `json:"whitelist"`
}

// CloudFrontSpecBehaviorWhitelist declares a whitelist of request parameters which are allowed.
type CloudFrontSpecBehaviorWhitelist struct {
	// Headers which will used when caching.
	Headers []string `json:"headers"`
	// Cookies which will be forwarded to the application.
	Cookies []string `json:"cookies"`
}

// CloudFrontSpecOrigin declares the origin which traffic will be sent.
type CloudFrontSpecOrigin struct {
	// Backend connection information for CloudFront.
	Endpoint string `json:"endpoint"`
	// Backend connection information for CloudFront.
	Policy string `json:"policy"`
	// How long CloudFront should wait before timing out.
	Timeout int64 `json:"timeout"`
}

// CloudFrontStatus defines the observed state of CloudFront
type CloudFrontStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Machine identifier for querying the CloudFront distribution.
	ID string `json:"arn,omitempty"`
	// Current state of the CloudFront distribution.
	State string `json:"state,omitempty"`
	// DomainName for creating CNAME records.
	DomainName string `json:"domainName,omitempty"`
	// How many invalidations are currently running.
	RunningInvalidations int64 `json:"runningInvalidations,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFront is the Schema for the cloudfronts API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CloudFront struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudFrontSpec   `json:"spec,omitempty"`
	Status CloudFrontStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFrontList contains a list of CloudFront
type CloudFrontList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudFront `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudFront{}, &CloudFrontList{})
}
