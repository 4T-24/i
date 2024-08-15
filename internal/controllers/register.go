package controllers

import (
	v1 "instancer/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *InstancierReconciler) RegisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		c.Spec.Slug = c.Name

		r.ctfdChallenges[c.Name] = &c.Spec
	case *v1.InstancedChallenge:
		r.challenges[c.Name] = c

		c.Spec.Slug = c.Name
		c.Spec.IsInstanced = true
		r.ctfdChallenges[c.Name] = &c.Spec.ChallengeSpec
	case *v1.OracleInstancedChallenge:
		r.challenges[c.Name] = c

		c.Spec.Slug = c.Name
		c.Spec.IsInstanced = true
		c.Spec.HasOracle = true
		r.ctfdChallenges[c.Name] = &c.Spec.ChallengeSpec
	}
}

func (r *InstancierReconciler) UnregisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		delete(r.ctfdChallenges, c.Name)
	case *v1.InstancedChallenge:
		delete(r.challenges, c.Name)
		delete(r.ctfdChallenges, c.Name)
	case *v1.OracleInstancedChallenge:
		delete(r.challenges, c.Name)
		delete(r.ctfdChallenges, c.Name)
	}
}
