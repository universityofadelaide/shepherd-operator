/*

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
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// Volumes which will be backed up.
	Volumes map[string]BackupSpecVolume `json:"volumes,omitempty"`
	// MySQL databases which will be backed up.
	MySQL map[string]BackupSpecMySQL `json:"mysql,omitempty"`
}

// BackupSpecVolume defines how to backup volumes.
type BackupSpecVolume struct {
	// ClaimName which will be backed up.
	ClaimName string `json:"claimName"`
}

// BackupSpecMySQL defines how to backup MySQL.
type BackupSpecMySQL struct {
	// Secret which will be used for connectivity.
	Secret BackupSpecMySQLSecret `json:"secret"`
}

type BackupSpecMySQLSecret struct {
	// Name of secret containing the mysql connection details.
	Name string `json:"name"`
	// Keys within secret to use for each parameter.
	Keys BackupSpecMySQLSecretKeys `json:"keys"`
}

// BackupSpecMySQLSecretKeys defines Secret keys for MySQL connectivity.
type BackupSpecMySQLSecretKeys struct {
	// Key which was applied to the application for database connectivity.
	Username string `json:"username"`
	// Key which was applied to the application for database connectivity.
	Password string `json:"password"`
	// Key which was applied to the application for database connectivity.
	Database string `json:"database"`
	// Key which was applied to the application for database connectivity.
	Hostname string `json:"hostname"`
	// Key which was applied to the application for database connectivity.
	Port string `json:"port"`
}

// BackupStatus defines the observed state of Backup
type BackupStatus struct {

}

// +kubebuilder:object:root=true

// Backup is the Schema for the backups API
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupList contains a list of Backup
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
