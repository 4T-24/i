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

// ChallengeHint defines the desired state of a hint
type ChallengeHint struct {
	Content      string                 `json:"content,omitempty"`
	Cost         int                    `json:"cost"`
	Requirements *ChallengeRequirements `json:"requirements,omitempty"`
}

// ChallengeRequirements defines the desired state of a requirement
type ChallengeRequirements struct {
	// Anonymize control the behavior of the resource if the prerequisites are
	// not validated:
	//  - if `nil`, defaults to `*false`
	//  - if `*false`, set the behavior as "hidden" (invisible until validated)
	//  - if `*true`, set the behavior to "anonymized" (visible but not much info)
	Anonymize *bool `json:"anonymize,omitempty"`

	// Prerequisites is the list of resources' ID that need to be validated in
	// order for the resource to meet its requirements.
	Prerequisites []int `json:"prerequisites"`
}

// ChallengeSpec defines the desired state of Challenge
type ChallengeSpec struct {
	// +kubebuilder:validation:Required
	// Name of the challenge
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// Name of the challenge
	Category string `json:"category"`

	// +kubebuilder:validation:Required
	// Description of the challenge
	Description string `json:"description"`

	// +kubebuilder:validation:Optional
	Value int `json:"value"`

	// +kubebuilder:validation:Optional
	Initial *int `json:"initial_value"`

	// +kubebuilder:validation:Optional
	Decay *int `json:"value_decay"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=linear;logarithmic
	DecayFunction string `json:"decay_function"`

	// +kubebuilder:validation:Optional
	Minimum *int `json:"minimum_value"`

	// +kubebuilder:validation:Optional
	MaxAttempts *int `json:"max_attempts"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=visible;hidden
	State string `json:"state"`

	// +kubebuilder:validation:Optional
	// Hints of the challenge
	Hints []ChallengeHint `json:"hints"`

	// +kubebuilder:validation:Optional
	// Requirements of the challenge
	Requirements ChallengeRequirements `json:"requirements"`

	// +kubebuilder:validation:Optional
	// Next challenge
	NextID *int `json:"next_id"`

	// +kubebuilder:validation:Required
	// Flag of this challenge
	Flag string `json:"flag"`

	// +kubebuilder:validation:Required
	// Type of the challenge
	Type string `json:"type"`
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
