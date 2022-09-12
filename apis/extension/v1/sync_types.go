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
	shpmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyncSpec defines the desired state of Sync
type SyncSpec struct {
	// Correlates to node ids in shepherd.
	Site       string     `json:"site"`
	BackupEnv  string     `json:"backupEnv"`
	RestoreEnv string     `json:"restoreEnv"`
	BackupSpec BackupSpec `json:"backupSpec"`
	// We can use the backup spec for the restore spec as we just need volumes/dbs.
	RestoreSpec BackupSpec `json:"restoreSpec"`
}

// SyncStatus defines the observed state of Sync
type SyncStatus struct {
	Backup  SyncStatusBackup  `json:"backup,omitempty"`
	Restore SyncStatusRestore `json:"restore,omitempty"`
}

// SyncStatusBackup defines the observed state of a Backup during a Sync.
type SyncStatusBackup struct {
	Name      string          `json:"name,omitempty"`
	Phase     shpmetav1.Phase `json:"phase"`
	StartTime *metav1.Time    `json:"startTime,omitempty"`
}

// SyncStatusRestore defines the observed state of a Restore during a Sync.
type SyncStatusRestore struct {
	Name           string          `json:"name,omitempty"`
	Phase          shpmetav1.Phase `json:"phase"`
	CompletionTime *metav1.Time    `json:"completionTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Start Time",type=string,JSONPath=.status.Backup.startTime
//+kubebuilder:printcolumn:name="Completion Time",type=string,JSONPath=.status.Restore.CompletionTime
//+kubebuilder:printcolumn:name="Backup Phase",type=string,JSONPath=.status.Backup.Phase
//+kubebuilder:printcolumn:name="Restore Phase",type=string,JSONPath=.status.Restore.Phase

// Sync is the Schema for the syncs API
type Sync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyncSpec   `json:"spec,omitempty"`
	Status SyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SyncList contains a list of Sync
type SyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Sync{}, &SyncList{})
}
