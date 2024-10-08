package templates

import (
	"fmt"
	v1 "instancer/api/v1"
	"strings"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceParams struct {
	Name         string
	Namespace    string
	CommonLabels map[string]string
	Ports        []v1.InstancedChallengePodPort
}

// apiVersion: v1
// kind: Service
// metadata:
//   name: {{ .Name }}
//   labels:
//     i.4ts.fr/pod: {{ .Name }}
//     {{- range $key, $value := .CommonLabels }}
//     {{ $key }}: "{{ $value }}"
//     {{- end }}
// spec:
//   selector:
//     i.4ts.fr/pod: {{ .Name }}
//     {{- range $key, $value := .CommonLabels }}
//     {{ $key }}: "{{ $value }}"
//     {{- end }}
//   ports:
//     {{- range .Ports }}
//     - name: "port-{{ . }}"
//       protocol: TCP
//       port: {{ . }}
//     {{- end }}

func NewService(p *ServiceParams) *core.Service {
	var labels = p.CommonLabels
	labels["i.4ts.fr/pod"] = p.Name

	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels:    labels,
		},
		Spec: core.ServiceSpec{
			Selector: labels,
		},
	}

	for _, port := range p.Ports {
		proto := core.ProtocolTCP
		if port.Protocol == "UDP" {
			proto = core.ProtocolUDP
		}

		service.Spec.Ports = append(service.Spec.Ports, core.ServicePort{
			Name:     strings.ToLower(fmt.Sprintf("port-%d-%s", port.Port, port.Protocol)),
			Protocol: proto,
			Port:     int32(port.Port),
		})
	}

	return service
}
