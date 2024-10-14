package controllers

import (
	"context"
	"instancer/internal/names"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *InstancierReconciler) DeleteInstance(challengeId, instanceId string) (*InstanceStatus, error) {
	var status = &InstanceStatus{}

	chall, found := r.GetChallengeSpec(challengeId)
	if !found {
		status.Status = "Unknown"
		return status, nil
	}

	status = &InstanceStatus{
		Name:    chall.Name,
		Timeout: chall.Timeout,
	}

	namespace := names.GetNamespaceName(challengeId, instanceId)

	var namespaceObj = &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespace,
		},
	}

	if err := r.Delete(context.Background(), namespaceObj); err != nil {
		return nil, err
	}

	delete(instancedIds, instanceId)
	status.Status = "Stopping"
	return status, nil
}
