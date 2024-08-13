package templates

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentParams holds the parameters for creating a Kubernetes Deployment.
type DeploymentParams struct {
	Name         string            // Name of the deployment
	Namespace    string            // Namespace of the deployment
	CommonLabels map[string]string // Common labels for the deployment
	Egress       string            // Value for the Egress label
	Spec         core.PodSpec      // Specification for the deployment's pod spec
}

// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       i.4ts.fr/pod: {{ .Name }}
//       {{- range $key, $value := .CommonLabels }}
//       {{ $key }}: "{{ $value }}"
//       {{- end }}
//   template:
//     metadata:
//       labels:
//         i.4ts.fr/pod: {{ .Name }}
//         i.4ts.fr/egress: "{{ .Egress }}"
//         {{- range $key, $value := .CommonLabels }}
//         {{ $key }}: "{{ $value }}"
//         {{- end }}
//     spec:
//       {{  asYaml .Spec 4 }}

func Optional[K any](k K) *K {
	return &k
}

func NewDeployment(p *DeploymentParams) *apps.Deployment {
	var labels = p.CommonLabels
	labels["i.4ts.fr/pod"] = p.Name

	deployment := &apps.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"io.kubernetes.cri-o.userns-mode": "auto:size=65536",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: Optional(int32(1)),
			Selector: &v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: labels,
				},
				Spec: p.Spec,
			},
		},
	}

	deployment.Spec.Template.Labels["i.4ts.fr/egress"] = p.Egress

	return deployment
}
