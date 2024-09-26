package templates

import (
	"fmt"
	"instancer/internal/env"

	emissary "github.com/emissary-ingress/emissary/v3/pkg/api/getambassador.io/v3alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IngressHost struct {
	Host        string
	ServiceName string
	ServicePort int
}

type IngressParams struct {
	Name      string
	Namespace string
	Kind      string
	Host      IngressHost
}

// apiVersion: voyager.appscode.com/v1
// kind: Ingress
// metadata:
//   name: {{ .Name }}
//   namespace: {{ .Namespace }}
//   annotations:
//     ingress.appscode.com/type: NodePort
//     ingress.appscode.com/use-node-port: "true"
// spec:
//   rules:
//   - host: {{ .Host }}
//     tcp:
//       nodePort: {{ .NodePort }}
//       port: {{ .ServicePort }}
//       backend:
//         service:
//           name: {{ .ServiceName }}
//           port:
//             number: {{ .ServicePort }}

func NewIngress(p *IngressParams) client.Object {
	if p.Kind == "tcp" {
		c := env.Get()

		var tcpmapping = &emissary.TCPMapping{
			ObjectMeta: v1.ObjectMeta{
				Name:      p.Name,
				Namespace: p.Namespace,
			},
			Spec: emissary.TCPMappingSpec{
				Port:    c.NodePort,
				Host:    p.Host.Host,
				Service: fmt.Sprintf("%s.%s:%d", p.Host.ServiceName, p.Namespace, p.Host.ServicePort),
			},
		}

		return tcpmapping
	}

	var ingress = &emissary.Mapping{
		ObjectMeta: v1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: emissary.MappingSpec{
			Prefix:   "/",
			Service:  fmt.Sprintf("%s.%s:%d", p.Host.ServiceName, p.Namespace, p.Host.ServicePort),
			Hostname: p.Host.Host,
			AllowUpgrade: []string{
				"websocket", // 🐝
			},
		},
	}

	return ingress
}
