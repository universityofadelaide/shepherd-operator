package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SolrSpec defines the desired state of Solr
type SolrSpec struct {
	// Name of the core which will be provisioned.
	Core string `json:"coreName"`
	// Version refers to the version/tag to use on the solr image repository.
	Version string `json:"version"`
	// Resources given to the Solr instance.
	Resources SolrSpecResources `json:"resources"`
}

// SolrSpecResources for provisioning Solr.
type SolrSpecResources struct {
	// CPU to use when provisioning Solr.
	CPU SolrSpecResourcesCPU `json:"cpu"`
	// Memory for the Solr instance. We don't set requests and limits because memory is important with Java applications.
	Memory resource.Quantity `json:"memory"`
	// Storage which is mounted for Solr.
	Storage resource.Quantity `json:"storage"`
}

// SolrSpecResourcesCPU for provisioning Solr.
type SolrSpecResourcesCPU struct {
	// CPU requests given to the Solr process.
	Request resource.Quantity `json:"request"`
	// CPU limits given to the Solr process.
	Limit resource.Quantity `json:"limit"`
}

// SolrStatus defines the observed state of Solr
type SolrStatus struct {
	// Host which applications will use to interact with Solr.
	Host string `json:"host"`
	// Port which applications will use to interact with Solr.
	Port int `json:"port"`
	// Core which applications will use for indexing.
	Core string `json:"core"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Solr is the Schema for the solrs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Solr struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SolrSpec   `json:"spec,omitempty"`
	Status SolrStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SolrList contains a list of Solr
type SolrList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Solr `json:"items"`
}

func init() {
	var object *Solr
	var list *SolrList
	SchemeBuilder.Register(object, list)
}
