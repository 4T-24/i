package controllers

import (
	"context"
	v1 "instancer/api/v1"
	"instancer/internal/names"
	"time"

	"codnect.io/chrono"
	"github.com/sirupsen/logrus"
)

type InstanceServers struct {
	Kind string `json:"kind"`
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
}

type InstanceStatus struct {
	Name    string     `json:"name"`
	Status  string     `json:"status"`
	Timeout int        `json:"timeout"`
	EndsAt  *time.Time `json:"endsAt,omitempty"`

	Servers []InstanceServers `json:"servers,omitempty"`
}

func (r *InstancierReconciler) GetChallengeSpec(challengeId string) (spec *v1.ChallengeSpec, found bool) {
	var chall *v1.Challenge
	var oraclechall *v1.OracleChallenge

	chall, found = r.availableChallenges[challengeId]
	if !found {
		oraclechall, found = r.availableOracleChallenges[challengeId]
		if found {
			spec = &oraclechall.Spec.ChallengeSpec
		}
	} else {
		spec = &chall.Spec
	}

	return
}

func (r *InstancierReconciler) GetInstance(challengeId, instanceId string) (*InstanceStatus, error) {
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

	namespace, err := r.GetNamespace(challengeId, instanceId)
	if err != nil {
		status.Status = "Stopped"
		return status, nil
	}

	if namespace.Status.Phase == "Terminating" {
		status.Status = "Stopping"
		return status, nil
	}

	for _, pod := range chall.Pods {
		deployment, err := r.GetDeployment(pod.Name, names.GetNamespaceName(challengeId, instanceId))
		if err != nil {
			status.Status = "Unknown"
			logrus.Error(err)
			return status, nil
		}
		if deployment.Status.AvailableReplicas == 0 {
			status.Status = "Starting"
			return status, nil
		}
	}

	status.Status = "Running"

	for _, port := range chall.ExposedPorts {
		status.Servers = append(status.Servers, r.GetServer(challengeId, instanceId, port.Pod, port.Kind))
	}

	t := namespace.CreationTimestamp.Time.Add(time.Duration(status.Timeout) * time.Second)
	status.EndsAt = &t

	namespaceName := names.GetNamespaceName(chall.Name, instanceId)
	if _, found := r.tasks[namespaceName]; !found {
		task, err := r.Schedule(func(ctx context.Context) {
			r.DeleteInstance(challengeId, instanceId)
		}, chrono.WithTime(time.Now().Add(time.Duration(status.Timeout)*time.Second)))
		if err == nil {
			r.tasks[namespaceName] = task
		}
	}

	return status, nil
}
