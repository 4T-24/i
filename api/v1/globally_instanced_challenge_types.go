/*
Copyright 2024.

Licensed under the BSD 3-Clause License
you may see the license in the LICENSE.md file
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GloballyInstancedChallengeOraclePodPort defines the desired state of GloballyInstancedChallengeOraclePodPort
type GloballyInstancedChallengeOraclePodPort struct {
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

// GloballyInstancedChallengeSpec defines the desired state of GlobalChallenge
type GloballyInstancedChallengeSpec struct {
	InstancedChallengeSpec `json:""`

	// +kubebuilder:validation:Required
	OraclePort GloballyInstancedChallengeOraclePodPort `json:"oraclePort"`
}

// +kubebuilder:object:root=true

// GloballyInstancedChallenge is the Schema for the GlobalChallenges API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category"
// +kubebuilder:printcolumn:name="Initial Value",type="integer",JSONPath=".spec.initial_value"
// +kubebuilder:printcolumn:name="Min Value",type="integer",JSONPath=".spec.minimum_value"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type GloballyInstancedChallenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GloballyInstancedChallengeSpec `json:"spec,omitempty"`
	Status ChallengeStatus                `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GloballyInstancedChallengeList contains a list of GlobalChallenge
type GloballyInstancedChallengeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GloballyInstancedChallenge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GloballyInstancedChallenge{}, &GloballyInstancedChallengeList{})
}
