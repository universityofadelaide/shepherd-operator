/*
Copyright 2022.

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

	shpmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// Volumes which will be backed up.
	Volumes map[string]shpmetav1.SpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]shpmetav1.SpecMySQL `json:"mysql,omitempty"`
}

// BackupStatus defines the observed state of Backup
type BackupStatus struct {
	StartTime      *metav1.Time    `json:"startTime,omitempty"`
	CompletionTime *metav1.Time    `json:"completionTime,omitempty"`
	ResticID       string          `json:"resticId,omitempty"`
	Phase          shpmetav1.Phase `json:"phase"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=.status.phase

// Backup is the Schema for the backups API
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BackupList contains a list of Backup
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
