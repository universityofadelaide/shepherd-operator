package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// IngressSpec defines the desired state of Ingress
type IngressSpec struct {
	// Rules which are used to Ingress traffic to an application.
	Routes IngressSpecRoutes `json:"routes"`
	// Whitelist rules for CloudFront.
	Whitelist awsv1beta1.CloudFrontSpecBehaviorWhitelist `json:"whitelist,omitempty"`
	// Backend connectivity details.
	Service    IngressSpecService    `json:"service"`
	Prometheus IngressSpecPrometheus `json:"prometheus"`
}

// IngressSpecRoutes declare the routes for the application.
type IngressSpecRoutes struct {
	// Primary domain and routing rule for the application.
	Primary IngressSpecRoute `json:"primary"`
	// Seconard domains and routing rules for the application.
	Secondary []IngressSpecRoute `json:"secondary"`
}

// IngressSpecRoute traffic from a domain and path to a service.
type IngressSpecRoute struct {
	// Domain used as part of a route rule.
	Domain string `json:"domain"`
	// Supaths included in the route rule.
	Subpaths []string `json:"subpaths"`
}

// IngressSpecService connects an Ingress to a Service.
type IngressSpecService struct {
	// Name of the Kubernetes Service object to route traffic to.
	Name string `json:"name"`
	// Port of the Kubernetes Service object to route traffic to.
	Port int `json:"port"`
}

// IngressSpecPrometheus defines the path which Prometheus can scrape application metrics.
type IngressSpecPrometheus struct {
	Path  string `json:"path"`
	Token string `json:"token"`
}

// IngressStatus defines the observed state of Ingress
type IngressStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Status of the provisioned CloudFront distribution.
	CloudFront IngressStatusCloudFrontRef
	// Status of the provisioned Certificate.
	Certificate IngressStatusCertificateRef
}

// IngressStatusCloudFrontRef provides status on the provisioned CloudFront.
type IngressStatusCloudFrontRef struct {
	// Name of the CloudFront distribution.
	Name string `json:"name,omitempty"`
	// Details on the provisioned CloudFront distribution.
	Details awsv1beta1.CloudFrontStatus
}

// IngressStatusCertificateRef provides status on the provisioned Certificate.
type IngressStatusCertificateRef struct {
	// Name of the certificate.
	Name string
	// Details on the provisioned certificate.
	Details awsv1beta1.CertificateStatus
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ingress is the Schema for the ingresses API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Ingress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngressSpec   `json:"spec,omitempty"`
	Status IngressStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressList contains a list of Ingress
type IngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ingress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ingress{}, &IngressList{})
}
