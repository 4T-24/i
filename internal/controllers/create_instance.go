package controllers

import (
	"context"
	"fmt"
	v1 "instancer/api/v1"
	"instancer/internal/names"
	"instancer/internal/templates"
	"strconv"
	"time"

	"codnect.io/chrono"
	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	instancedIds       = make(map[string]string)
	instanceInCreation = make(map[string]bool)

	retryPolicy = retrypolicy.Builder[any]().
			WithDelay(time.Second).
			WithMaxRetries(3).
			Build()
)

func (r *InstancierReconciler) CreateInstance(challengeId, instanceId string) (*InstanceStatus, error) {
	var status = &InstanceStatus{}

	chall, found := r.GetChallengeSpec(challengeId)
	if !found {
		status.Status = "Unknown"
		return status, nil
	}

	if v, found := instancedIds[instanceId]; found {
		// Delete the instance if it already exists
		r.DeleteInstance(v, instanceId)
		delete(instancedIds, instanceId)
	}

	if _, found := instanceInCreation[instanceId]; found {
		status.Status = "Creating"
		return status, nil
	}

	instanceInCreation[instanceId] = true
	defer func() {
		// Avoid locking the instance permanently
		delete(instanceInCreation, instanceId)
	}()

	status = &InstanceStatus{
		Name:    chall.Name,
		Timeout: chall.Timeout,
	}

	id := names.GetId()
	namespace := names.GetNamespaceName(challengeId, instanceId)
	commonLabels := names.GetCommonLabels(challengeId, instanceId, id)

	var namespaceObj = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: commonLabels,
		},
	}
	namespaceObj.Labels["i.4ts.fr/ttl"] = fmt.Sprint(chall.Timeout)
	namespaceObj.Labels["i.4ts.fr/stops-at-timestamp"] = fmt.Sprint(time.Now().Add(time.Duration(chall.Timeout) * time.Second).Unix())

	err := failsafe.Run(func() error {
		return r.Create(context.Background(), namespaceObj)
	}, retryPolicy)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task, err := r.Schedule(func(ctx context.Context) {
		r.DeleteInstance(challengeId, instanceId)
	}, chrono.WithTime(time.Now().Add(time.Duration(status.Timeout)*time.Second)))
	if err == nil {
		r.tasks[namespace] = task
	}

	if chall.RegistrySecret != nil {
		var secret corev1.Secret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      chall.RegistrySecret.Name,
				Namespace: chall.RegistrySecret.Namespace,
			},
		}
		if err := r.Get(context.Background(), client.ObjectKeyFromObject(&secret), &secret); err != nil {
			logrus.Error(err)
			return nil, err
		}

		var newSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      chall.RegistrySecret.Name,
				Namespace: namespace,
			},
			Data: secret.Data,
			Type: secret.Type,
		}

		if err := r.Create(context.Background(), &newSecret); err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	networkpolicies := templates.NewNetworkPolicy(&templates.NetworkPolicyParams{
		Namespace:    namespace,
		CommonLabels: commonLabels,
		Pods:         chall.Pods,
	})

	for _, networkpolicy := range networkpolicies {
		err := failsafe.Run(func() error {
			return r.Create(context.Background(), networkpolicy)
		}, retryPolicy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	for _, pod := range chall.Pods {
		annotations := make(map[string]string)
		for _, p := range pod.Ports {
			annotations[fmt.Sprintf("i.4ts.fr/hostname-%d", p.Port)] = names.GetHost(pod.Name, p.Port, challengeId, id)
			annotations[fmt.Sprintf("i.4ts.fr/uri-%d", p.Port)] = "https://" + names.GetHost(pod.Name, p.Port, challengeId, id)
		}

		deployment := templates.NewDeployment(&templates.DeploymentParams{
			Name:         pod.Name,
			Namespace:    namespace,
			CommonLabels: commonLabels,
			Egress:       strconv.FormatBool(pod.Egress),
			Spec:         pod.Spec,
			Annotations:  annotations,
		})
		err := failsafe.Run(func() error {
			return r.Create(context.Background(), deployment)
		}, retryPolicy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		service := templates.NewService(&templates.ServiceParams{
			Name:         pod.Name,
			Namespace:    namespace,
			CommonLabels: commonLabels,
			Ports:        pod.Ports,
		})
		err = failsafe.Run(func() error {
			return r.Create(context.Background(), service)
		}, retryPolicy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	for i, port := range chall.ExposedPorts {
		ingress := templates.NewIngress(&templates.IngressParams{
			Name:      fmt.Sprintf("%s-%d", port.Pod, i),
			Namespace: namespace,
			Kind:      port.Kind,
			Host: templates.IngressHost{
				Host:        names.GetHost(port.Pod, port.Port, challengeId, id),
				ServiceName: port.Pod,
				ServicePort: port.Port,
			},
		})
		err := failsafe.Run(func() error {
			return r.Create(context.Background(), ingress)
		}, retryPolicy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	instancedIds[instanceId] = challengeId

	status.Status = "Starting"
	return status, nil
}

func (r *InstancierReconciler) CreateGlobalInstances() ([]*InstanceStatus, error) {
	var out []*InstanceStatus

	for _, chall := range r.challenges {
		switch v := chall.(type) {
		case *v1.GloballyInstancedChallenge:
			status, err := r.CreateInstance(v.Name, "global")
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
