package controllers

import (
	"context"
	"instancer/internal/env"
	"instancer/internal/names"

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

func (r *InstancierReconciler) GetServer(challengeId, instanceId, pod, kind string) InstanceServers {
	s := InstanceServers{
		Kind: kind,
		Host: names.GetHost(pod, challengeId, instanceId),
	}

	if kind == "tcp" {
		c := env.Get()
		s.Port = c.NodePort
	}

	return s
}
