package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
)

// ImageSpec defines the desired state of Image
type ImageSpec struct {
	// Connection details for the mysql image.
	Connection Connection `json:"connection"`
	// Rules for the mysql image.
	Rules ImageSpecRules `json:"rules"`
	// Destinations to push the  mysql image to.
	Destinations []string `json:"destinations"`
}

// Connection defines how to connect to MySQL.
type Connection struct {
	// ConfigMap which will be used for connectivity.
	ConfigMap ConnectionConfigMap `json:"configmap"`
	// Secret which will be used for connectivity.
	Secret ConnectionSecret `json:"secret"`
}

// ConnectionConfigMap defines connection information for MySQL.
type ConnectionConfigMap struct {
	// Name of the ConfigMap.
	Name string `json:"name"`
	// ConfigMap keys used for connectivity.
	Keys ConnectionConfigMapKeys `json:"keys"`
}

// ConnectionConfigMapKeys defines ConfigMap keys for MySQL connectivity.
type ConnectionConfigMapKeys struct {
	// Database which was applied to the application for database connectivity.
	Database string `json:"database"`
	// Hostname which was applied to the application for database connectivity.
	Hostname string `json:"hostname"`
	// Port which was applied to the application for database connectivity.
	Port string `json:"port"`
}

// ConnectionSecret defines connection information for MySQL.
type ConnectionSecret struct {
	// Name of the Secret.
	Name string `json:"name"`
	// Secret keys used for connectivity.
	Keys ConnectionSecretKeys `json:"keys"`
}

// ConnectionSecretKeys defines Secret keys for MySQL connectivity.
type ConnectionSecretKeys struct {
	// Username which was applied to the application for database connectivity.
	Username string `json:"username"`
	// Password which was applied to the application for database connectivity.
	Password string `json:"password"`
}

// ImageSpecRules for configuring dump.
// This has been lifted from the github.com/skpr/mtk/dump project.
// We have copied the configuration into this API package so Kubebuilder can
// generate the required DeepCopy() functions.
type ImageSpecRules struct {
	Rewrite map[string]ImageSpecRulesRewrite `yaml:"rewrite" json:"rewrite"`
	NoData  []string                         `yaml:"nodata"  json:"nodata"`
	Ignore  []string                         `yaml:"ignore"  json:"ignore"`
}

// ImageSpecRulesRewrite rules for while dumping a database.
type ImageSpecRulesRewrite map[string]string

// ImageStatus defines the observed state of Image
type ImageStatus struct {
	// Current state of the database image.
	Phase skprmetav1.Phase `json:"phase,omitempty"`
	// Time which the MySQL image started building
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Time which the MySQL image finished building
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Image is the Schema for the images API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageList contains a list of Image
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
