package templates

import (
	"instancer/internal/utils"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentParams holds the parameters for creating a Kubernetes Deployment.
type DeploymentParams struct {
	Name         string            // Name of the deployment
	Namespace    string            // Namespace of the deployment
	Egress       string            // Value for the Egress label
	Annotations  map[string]string // Annotations for the deployment
	CommonLabels map[string]string // Common labels for the deployment
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

func NewDeployment(p *DeploymentParams) *apps.Deployment {
	var labels = p.CommonLabels
	labels["i.4ts.fr/pod"] = p.Name

	p.Annotations["io.kubernetes.cri-o.userns-mode"] = "auto:size=65536"

	deployment := &apps.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:        p.Name,
			Namespace:   p.Namespace,
			Labels:      labels,
			Annotations: p.Annotations,
		},
		Spec: apps.DeploymentSpec{
			Replicas: utils.Optional(int32(1)),
			Selector: &v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels:      labels,
					Annotations: p.Annotations,
				},
				Spec: p.Spec,
			},
		},
	}

	deployment.Spec.Template.Labels["i.4ts.fr/egress"] = p.Egress

	return deployment
}
