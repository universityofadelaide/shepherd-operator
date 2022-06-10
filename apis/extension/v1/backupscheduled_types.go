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

// BackupScheduledSpec defines the desired state of BackupScheduled
type BackupScheduledSpec struct {
	// Retention defines how backup retention behavior.
	Retention shpmetav1.RetentionSpec `json:"retention"`
	// Schedule is the crontab statement which defines how often a backup should run.
	Schedule shpmetav1.ScheduledSpec `json:"schedule"`
	// Volumes which will be backed up.
	Volumes map[string]shpmetav1.SpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]shpmetav1.SpecMySQL `json:"mysql,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BackupScheduled is the Schema for the backupscheduleds API
type BackupScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduledSpec       `json:"spec,omitempty"`
	Status shpmetav1.ScheduledStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BackupScheduledList contains a list of BackupScheduled
type BackupScheduledList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupScheduled `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupScheduled{}, &BackupScheduledList{})
}
