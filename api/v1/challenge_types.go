/*
Copyright 2024.

Licensed under the BSD 3-Clause License
you may see the license in the LICENSE.md file
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChallengeHint defines the desired state of a hint
type ChallengeHint struct {
	Content      string            `json:"content,omitempty"`
	Cost         int               `json:"cost"`
	Requirements *HintRequirements `json:"requirements,omitempty"`
}

// ChallengeFile defines the desired state of a files
type ChallengeFile struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Path string `json:"path"`
}

// HintRequirements defines the desired state of a requirement
type HintRequirements struct {
	// Anonymize control the behavior of the resource if the prerequisites are
	// not validated:
	//  - if `nil`, defaults to `*false`
	//  - if `*false`, set the behavior as "hidden" (invisible until validated)
	//  - if `*true`, set the behavior to "anonymized" (visible but not much info)
	Anonymize *bool `json:"anonymize,omitempty"`

	// Prerequisites is the list of resources' slug or name that need to be validated in
	// order for the resource to meet its requirements.
	Prerequisites []int `json:"prerequisites"`
}

// ChallengeRequirements defines the desired state of a requirement
type ChallengeRequirements struct {
	// Anonymize control the behavior of the resource if the prerequisites are
	// not validated:
	//  - if `nil`, defaults to `*false`
	//  - if `*false`, set the behavior as "hidden" (invisible until validated)
	//  - if `*true`, set the behavior to "anonymized" (visible but not much info)
	Anonymize *bool `json:"anonymize,omitempty"`

	// Prerequisites is the list of resources' slug or name that need to be validated in
	// order for the resource to meet its requirements.
	Prerequisites []string `json:"prerequisites"`
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
	// Files of the challenge
	Files []ChallengeFile `json:"files"`

	// +kubebuilder:validation:Optional
	// Repository to get the files
	Repository string `json:"repository"`

	// +kubebuilder:validation:Optional
	// Requirements of the challenge
	Requirements ChallengeRequirements `json:"requirements"`

	// +kubebuilder:validation:Optional
	// Next challenge
	NextSlug string `json:"next_slug"`

	// +kubebuilder:validation:Required
	// Flag of this challenge
	Flag string `json:"flag"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=i_static;i_dynamic
	// Type of the challenge
	Type string `json:"type"`

	// Field for later use
	Slug        string `json:"-"`
	IsInstanced bool   `json:"-"`
	IsGlobal    bool   `json:"-"`
	HasOracle   bool   `json:"-"`
}

// ChallengeStatus defines the observed state of Challenge
type ChallengeStatus struct {
	Phase string `json:"phase"`
	Error string `json:"error"`
}

// Challenge is the Schema for the challenges API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category"
// +kubebuilder:printcolumn:name="Initial Value",type="integer",JSONPath=".spec.initial_value"
// +kubebuilder:printcolumn:name="Min Value",type="integer",JSONPath=".spec.minimum_value"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Challenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChallengeSpec   `json:"spec,omitempty"`
	Status ChallengeStatus `json:"status,omitempty"`
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
