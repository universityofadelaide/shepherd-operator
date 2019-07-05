package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	// Provisioner used to create databases.
	Provisioner string `json:"provisioner"`
	// Privileges which the application requires.
	Privileges []string `json:"privileges"`
}

// Phase which indicates the status of an object.
type Phase string

const (
	// PhaseFailed to be assigned when the database provisioning fails.
	PhaseFailed Phase = "Failed"
	// PhaseReady to be assigned when the database is ready for connections.
	PhaseReady Phase = "Ready"
)

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	// Used for determining if an APIs information is up to date.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Current state of the database being provisioned.
	Phase Phase `json:"phase"`
	// Connection details for the database.
	Connection DatabaseStatusConnection `json:"connection"`
}

// DatabaseStatusConnection for applications.
type DatabaseStatusConnection struct {
	// Hostname used when connecting to the database.
	Hostname string `json:"hostname,omitempty"`
	// Port used when connecting to the database.
	Port int `json:"port,omitempty"`
	// Database used when connecting to the database.
	Database string `json:"database,omitempty"`
	// Username used when connecting to the database.
	Username string `json:"username,omitempty"`
	// Password used when connecting to the database.
	Password string `json:"password,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Database is the Schema for the databaseclaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
