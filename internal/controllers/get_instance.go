package controllers

import (
	"errors"
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/names"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type InstanceServers struct {
	Kind         string `json:"kind"`
	Host         string `json:"host"`
	Port         int    `json:"port,omitempty"`
	Description  string `json:"description,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type InstanceStatus struct {
	Name    string     `json:"name"`
	Status  string     `json:"status"`
	Timeout int        `json:"timeout"`
	EndsAt  *time.Time `json:"endsAt,omitempty"`

	Servers []InstanceServers `json:"servers,omitempty"`
}

func (r *InstancierReconciler) GetChallengeSpec(challengeId string) (*v1.InstancedChallengeSpec, bool) {
	obj, found := r.challenges[challengeId]
	if !found {
		return nil, false
	}

	switch o := obj.(type) {
	case *v1.InstancedChallenge:
		return &o.Spec, true
	case *v1.GloballyInstancedChallenge:
		return &o.Spec.InstancedChallengeSpec, true
	case *v1.OracleInstancedChallenge:
		return &o.Spec.InstancedChallengeSpec, true
	default:
		return nil, false
	}
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

	id := namespace.Labels["i.4ts.fr/instance"]

	for _, port := range chall.ExposedPorts {
		status.Servers = append(status.Servers, r.GetServer(challengeId, id, port))
	}

	t := namespace.CreationTimestamp.Time.Add(time.Duration(status.Timeout) * time.Second)
	status.EndsAt = &t

	return status, nil
}

func (r *InstancierReconciler) GetGlobalInstances() ([]*InstanceStatus, error) {
	var out []*InstanceStatus
	for _, chall := range r.challenges {
		switch v := chall.(type) {
		case *v1.GloballyInstancedChallenge:
			status, err := r.GetInstance(v.Name, "global")
			if err != nil {
				out = append(out, &InstanceStatus{
					Name:   v.Name,
					Status: "Errored : " + err.Error(),
				})
				continue
			}
			out = append(out, status)
		}
	}
	return out, nil
}

func (r *InstancierReconciler) IsInstanceSolved(challengeId, instanceId string) (bool, error) {
	var err error

	chall, found := r.challenges[challengeId]
	if !found {
		return false, errors.New("challenge is not an oracle challenge")
	}

	var oraclechall *v1.OracleInstancedChallenge

	switch v := chall.(type) {
	case *v1.OracleInstancedChallenge:
		oraclechall = v
	default:
		return false, err
	}

	namespace := names.GetNamespaceName(challengeId, instanceId)
	uri := fmt.Sprintf("http://%s.%s:%d/", oraclechall.Spec.OraclePort.Pod, namespace, oraclechall.Spec.OraclePort.Port)
	uri, err = url.JoinPath(uri, oraclechall.Spec.OraclePort.Route)
	if err != nil {
		return false, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			logrus.Info("Cannot reach host, am I on kubernetes ?")
			return false, nil
		}
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
