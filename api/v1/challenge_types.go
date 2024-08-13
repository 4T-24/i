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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChallengePodPort defines the desired state of ChallengePodPort
type ChallengePodPort struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Required
	Port int `json:"port"`

	// +kubebuilder:validation:Enum=TCP;UDP
	Protocol string `json:"protocol"`
}

// ChallengeExposedPort defines the desired state of ChallengeExposedPort
type ChallengeExposedPort struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Required
	Port int `json:"port"`

	// +kubebuilder:validation:Enum=tcp;http
	Kind string `json:"kind"`

	// +kubebuilder:validation:Required
	Pod string `json:"pod"`
}

// ChallengePod defines the desired state of ChallengePod
type ChallengePod struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Ports []ChallengePodPort `json:"ports"`

	Egress bool `json:"egress"`

	// +kubebuilder:validation:Required
	Spec v1.PodSpec `json:"spec"`
}

// ChallengeSpec defines the desired state of Challenge
type ChallengeSpec struct {
	// +kubebuilder:validation:Required
	// Name of the challenge
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// Timeout in seconds, after which the challenge is deleted
	Timeout int `json:"timeout"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	ExposedPorts []ChallengeExposedPort `json:"exposedPorts"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// Pods to deploy for the challenge
	Pods []ChallengePod `json:"pods"`
}

// +kubebuilder:object:root=true

// Challenge is the Schema for the challenges API
type Challenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ChallengeSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ChallengeList contains a list of Challenge
type ChallengeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Challenge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Challenge{}, &ChallengeList{})
}
