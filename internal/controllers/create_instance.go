package controllers

import (
	"context"
	"fmt"
	"instancer/internal/names"
	"instancer/internal/templates"
	"strconv"
	"time"

	"codnect.io/chrono"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *InstancierReconciler) CreateInstance(challengeId, instanceId string) (*InstanceStatus, error) {
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

	id := names.GetId()
	namespace := names.GetNamespaceName(challengeId, instanceId)
	commonLabels := names.GetCommonLabels(challengeId, instanceId, id)

	var namespaceObj = &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:   namespace,
			Labels: commonLabels,
		},
	}
	namespaceObj.Labels["i.4ts.fr/ttl"] = fmt.Sprint(chall.Timeout)

	if err := r.Create(context.Background(), namespaceObj); err != nil {
		return nil, err
	}

	task, err := r.Schedule(func(ctx context.Context) {
		r.DeleteInstance(challengeId, instanceId)
	}, chrono.WithTime(time.Now().Add(time.Duration(status.Timeout)*time.Second)))
	if err == nil {
		r.tasks[namespace] = task
	}

	networkpolicies := templates.NewNetworkPolicy(&templates.NetworkPolicyParams{
		Namespace:    namespace,
		CommonLabels: commonLabels,
		Pods:         chall.Pods,
	})

	for _, networkpolicy := range networkpolicies {
		if err := r.Create(context.Background(), networkpolicy); err != nil {
			return nil, err
		}
	}

	for _, pod := range chall.Pods {
		deployment := templates.NewDeployment(&templates.DeploymentParams{
			Name:         pod.Name,
			Namespace:    namespace,
			CommonLabels: commonLabels,
			Egress:       strconv.FormatBool(pod.Egress),
			Spec:         pod.Spec,
		})
		if err := r.Create(context.Background(), deployment); err != nil {
			return nil, err
		}

		service := templates.NewService(&templates.ServiceParams{
			Name:         pod.Name,
			Namespace:    namespace,
			CommonLabels: commonLabels,
			Ports:        pod.Ports,
		})
		if err := r.Create(context.Background(), service); err != nil {
			return nil, err
		}
	}

	for i, port := range chall.ExposedPorts {
		ingress := templates.NewIngress(&templates.IngressParams{
			Name:      fmt.Sprintf("%s-%d", port.Pod, i),
			Namespace: namespace,
			Kind:      port.Kind,
			Host: templates.IngressHost{
				Host:        names.GetHost(port.Pod, challengeId, id),
				ServiceName: port.Pod,
				ServicePort: port.Port,
			},
		})
		if err := r.Create(context.Background(), ingress); err != nil {
			return nil, err
		}
	}

	status.Status = "Starting"
	return status, nil
}
