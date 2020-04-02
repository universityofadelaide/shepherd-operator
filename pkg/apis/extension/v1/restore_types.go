/*
Copyright 2019 University of Adelaide.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	shpmetav1 "github.com/universityofadelaide/shepherd-operator/pkg/apis/meta/v1"
)

// RestoreSpec defines the desired state of Restore
type RestoreSpec struct {
	// Name of the backup to restore from.
	BackupName string `json:"backupName"`
	// Volumes which will be backed up.
	Volumes map[string]shpmetav1.SpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]shpmetav1.SpecMySQL `json:"mysql,omitempty"`
}

// RestoreStatus defines the observed state of Restore
type RestoreStatus struct {
	StartTime      *metav1.Time    `json:"startTime,omitempty"`
	CompletionTime *metav1.Time    `json:"completionTime,omitempty"`
	Phase          shpmetav1.Phase `json:"phase"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Restore is the Schema for the restores API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=.status.phase
type Restore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestoreSpec   `json:"spec,omitempty"`
	Status RestoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RestoreList contains a list of Restore
type RestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Restore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Restore{}, &RestoreList{})
}
