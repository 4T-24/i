package controllers

import (
	v1 "instancer/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *InstancierReconciler) RegisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		c.Spec.Slug = c.Name

		r.challenges[c.Name] = c
		r.ctfdChallengesSpecs[c.Name] = &c.Spec
	case *v1.InstancedChallenge:
		c.Spec.ChallengeSpec.Slug = c.Name
		c.Spec.ChallengeSpec.IsInstanced = true

		r.challenges[c.Name] = c
		r.ctfdChallengesSpecs[c.Name] = &c.Spec.ChallengeSpec
	case *v1.GloballyInstancedChallenge:
		c.Spec.ChallengeSpec.Slug = c.Name
		c.Spec.ChallengeSpec.IsGlobal = true

		r.challenges[c.Name] = c
		r.ctfdChallengesSpecs[c.Name] = &c.Spec.ChallengeSpec
	case *v1.OracleInstancedChallenge:
		c.Spec.ChallengeSpec.Slug = c.Name
		c.Spec.ChallengeSpec.IsInstanced = true
		c.Spec.ChallengeSpec.HasOracle = true

		r.challenges[c.Name] = c
		r.ctfdChallengesSpecs[c.Name] = &c.Spec.ChallengeSpec
	}
}

func (r *InstancierReconciler) UnregisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		delete(r.ctfdChallengesSpecs, c.Name)
	case *v1.InstancedChallenge:
		delete(r.challenges, c.Name)
		delete(r.ctfdChallengesSpecs, c.Name)
	case *v1.OracleInstancedChallenge:
		delete(r.challenges, c.Name)
		delete(r.ctfdChallengesSpecs, c.Name)
	}
}
