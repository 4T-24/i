/*
Copyright 2024.

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

// ChallengeOraclePodPort defines the desired state of ChallengeOraclePodPort
type ChallengeOraclePodPort struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Required
	Port int `json:"port"`

	// +kubebuilder:validation:Required
	Pod string `json:"pod"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default=/is_solved
	Route string `json:"route"`
}

// OracleChallengeSpec defines the desired state of OracleChallenge
type OracleChallengeSpec struct {
	ChallengeSpec `json:""`

	// +kubebuilder:validation:Required
	OraclePort ChallengeOraclePodPort `json:"oraclePort"`
}

// +kubebuilder:object:root=true

// OracleChallenge is the Schema for the oraclechallenges API
type OracleChallenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OracleChallengeSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// OracleChallengeList contains a list of OracleChallenge
type OracleChallengeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OracleChallenge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OracleChallenge{}, &OracleChallengeList{})
}
