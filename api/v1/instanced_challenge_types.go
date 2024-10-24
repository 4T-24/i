/*
Copyright 2024.

Licensed under the BSD 3-Clause License
you may see the license in the LICENSE.md file
*/

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstancedChallengePodPort defines the desired state of InstancedChallengePodPort
type InstancedChallengePodPort struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Required
	Port int `json:"port"`

	// +kubebuilder:validation:Enum=TCP;UDP
	Protocol string `json:"protocol"`
}

// InstancedChallengeExposedPort defines the desired state of InstancedChallengeExposedPort
type InstancedChallengeExposedPort struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Required
	Port int `json:"port"`

	// +kubebuilder:validation:Enum=tcp;http
	Kind string `json:"kind"`

	// +kubebuilder:validation:Required
	Pod string `json:"pod"`

	// +kubebuilder:validation:Optional
	Description string `json:"description"`
}

// InstancedChallengePod defines the desired state of InstancedChallengePod
type InstancedChallengePod struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Ports []InstancedChallengePodPort `json:"ports"`

	Egress bool `json:"egress"`

	// +kubebuilder:validation:Required
	Spec v1.PodSpec `json:"spec"`
}

// InstancedChallengeRegistrySecret defines the desired state of InstancedChallengeRegistrySecret
type InstancedChallengeRegistrySecret struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
}

// InstancedChallengeSpec defines the desired state of Challenge
type InstancedChallengeSpec struct {
	ChallengeSpec `json:""`

	// +kubebuilder:validation:Required
	// Timeout in seconds, after which the challenge is deleted
	Timeout int `json:"timeout"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	ExposedPorts []InstancedChallengeExposedPort `json:"exposedPorts"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// Pods to deploy for the challenge
	Pods []InstancedChallengePod `json:"pods"`

	// +kubebuilder:validation:Optional
	// Registry secret to use for pulling images
	RegistrySecret *InstancedChallengeRegistrySecret `json:"registrySecret"`
}

// +kubebuilder:object:root=true

// InstancedChallenge is the Schema for the challenges API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category"
// +kubebuilder:printcolumn:name="Initial Value",type="integer",JSONPath=".spec.initial_value"
// +kubebuilder:printcolumn:name="Min Value",type="integer",JSONPath=".spec.minimum_value"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type InstancedChallenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstancedChallengeSpec `json:"spec,omitempty"`
	Status ChallengeStatus        `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstancedChallengeList contains a list of Challenge
type InstancedChallengeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstancedChallenge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InstancedChallenge{}, &InstancedChallengeList{})
}
