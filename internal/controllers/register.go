package controllers

import (
	v1 "instancer/api/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

// SetupWithManager sets up the controller with the Manager.
func (r *InstancierReconciler) RegisterChallenge(obj runtime.Object) {
	switch c := obj.(type) {
	case *v1.Challenge:
		r.availableChallenges[c.Name] = c
	case *v1.OracleChallenge:
		r.availableOracleChallenges[c.Name] = c
	}
}
