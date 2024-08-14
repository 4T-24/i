package controllers

import (
	v1 "instancer/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *InstancierReconciler) RegisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		r.ctfdChallenges = append(r.ctfdChallenges, &c.Spec)
	case *v1.InstancedChallenge:
		r.challenges[c.Name] = c
		r.ctfdChallenges = append(r.ctfdChallenges, &c.Spec.ChallengeSpec)
	case *v1.OracleInstancedChallenge:
		r.challenges[c.Name] = c
		r.ctfdChallenges = append(r.ctfdChallenges, &c.Spec.ChallengeSpec)
	}
}

func (r *InstancierReconciler) UnregisterChallenge(obj client.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		for i, challenge := range r.ctfdChallenges {
			if challenge.Name == c.Spec.Name {
				r.ctfdChallenges = append(r.ctfdChallenges[:i], r.ctfdChallenges[i+1:]...)
				break
			}
		}
	case *v1.InstancedChallenge:
		delete(r.challenges, c.Name)
		for i, challenge := range r.ctfdChallenges {
			if challenge.Name == c.Spec.Name {
				r.ctfdChallenges = append(r.ctfdChallenges[:i], r.ctfdChallenges[i+1:]...)
				break
			}
		}
	case *v1.OracleInstancedChallenge:
		delete(r.challenges, c.Name)
		for i, challenge := range r.ctfdChallenges {
			if challenge.Name == c.Spec.Name {
				r.ctfdChallenges = append(r.ctfdChallenges[:i], r.ctfdChallenges[i+1:]...)
				break
			}
		}
	}
}
