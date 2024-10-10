package controllers

import (
	"context"
	"fmt"
	"instancer/internal/env"
	"instancer/internal/names"

	instancer "instancer/api/v1"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *InstancierReconciler) GetNamespace(challengeId, instanceId string) (*corev1.Namespace, error) {
	var namespace = corev1.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: names.GetNamespaceName(challengeId, instanceId),
		},
	}
	err := r.Get(context.Background(), client.ObjectKeyFromObject(&namespace), &namespace)
	return &namespace, err
}

func (r *InstancierReconciler) GetDeployment(pod, namespace string) (*v1.Deployment, error) {
	var deployment = &v1.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      pod,
			Namespace: namespace,
		},
	}
	err := r.Get(context.Background(), client.ObjectKeyFromObject(deployment), deployment)
	return deployment, err
}

func (r *InstancierReconciler) GetServer(challengeId string, id string, port instancer.InstancedChallengeExposedPort) InstanceServers {
	s := InstanceServers{
		Kind:        port.Kind,
		Host:        names.GetHost(port.Pod, port.Port, challengeId, id),
		Description: port.Description,
	}

	if port.Kind == "tcp" {
		c := env.Get()
		s.Port = c.NodePort
		s.Instructions = fmt.Sprintf("openssl s_client -quiet -verify_quiet -connect %s:%d", s.Host, s.Port)
	}

	return s
}
