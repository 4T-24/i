package templates

import (
	v1 "instancer/api/v1"
	"instancer/internal/utils"

	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// apiVersion: networking.k8s.io/v1
// kind: NetworkPolicy
// metadata:
//   name: isolate-network
// spec:
//   podSelector: {}
//   policyTypes:
//     - Ingress
//     - Egress
//   ingress:
//     - from:
//         - namespaceSelector:
//             matchLabels:
//               {{- range $key, $value := .CommonLabels }}
//               {{ $key }}: "{{ $value }}"
//               {{- end }}
//   egress:
//     - to:
//         - namespaceSelector:
//             matchLabels:
//               {{- range $key, $value := .CommonLabels }}
//               {{ $key }}: "{{ $value }}"
//               {{- end }}
//     - to:
//         - namespaceSelector:
//             matchLabels:
//               kubernetes.io/metadata.name: kube-system
//       ports:
//         - protocol: UDP
//           port: 53
// ---
// apiVersion: networking.k8s.io/v1
// kind: NetworkPolicy
// metadata:
//   name: allow-ingress
// spec:
//   podSelector:
//     matchLabels:
//       i.4ts.fr/pod: {{ .Pod }}
//   policyTypes:
//     - Ingress
//   ingress:
//     - from:
//       - namespaceSelector:
//           matchLabels:
//             kubernetes.io/metadata.name: emissary
//       - podSelector:
//           matchLabels:
//              app.kubernetes.io/instance: emissary-ingress
// ---
// apiVersion: networking.k8s.io/v1
// kind: NetworkPolicy
// metadata:
//   name: allow-egress
// spec:
//   podSelector:
//     matchLabels:
//       i.4ts.fr/egress: "true"
//   policyTypes:
//     - Egress
//   egress:
//     - to:
//         - ipBlock:
//             cidr: "0.0.0.0/0"
//             except:
//               - "10.0.0.0/8"
//               - "192.168.0.0/16"
//               - "172.16.0.0/20"

// NetworkPolicyParams holds the parameters for creating a Kubernetes NetworkPolicy.
type NetworkPolicyParams struct {
	Namespace    string            // Namespace of the NetworkPolicy
	CommonLabels map[string]string // Common labels for namespace selectors
	Pods         []v1.InstancedChallengePod
}

func NewNetworkPolicy(p *NetworkPolicyParams) []*networking.NetworkPolicy {
	var networkPolicies []*networking.NetworkPolicy

	networkPolicies = append(networkPolicies, &networking.NetworkPolicy{
		ObjectMeta: meta.ObjectMeta{
			Name:      "isolate-network",
			Namespace: p.Namespace,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: meta.LabelSelector{},
			PolicyTypes: []networking.PolicyType{
				networking.PolicyTypeIngress,
				networking.PolicyTypeEgress,
			},
			Ingress: []networking.NetworkPolicyIngressRule{
				{
					From: []networking.NetworkPolicyPeer{
						{
							NamespaceSelector: &meta.LabelSelector{
								MatchLabels: p.CommonLabels,
							},
						},
					},
				},
			},
			Egress: []networking.NetworkPolicyEgressRule{
				{
					To: []networking.NetworkPolicyPeer{
						{
							NamespaceSelector: &meta.LabelSelector{
								MatchLabels: p.CommonLabels,
							},
						},
					},
				},
				{
					To: []networking.NetworkPolicyPeer{
						{
							NamespaceSelector: &meta.LabelSelector{
								MatchLabels: map[string]string{
									"kubernetes.io/metadata.name": "kube-system",
								},
							},
						},
					},
					Ports: []networking.NetworkPolicyPort{
						{
							Protocol: utils.Optional(core.ProtocolUDP),
							Port:     utils.Optional(intstr.FromInt32(53)),
						},
					},
				},
			},
		},
	})

	for _, pod := range p.Pods {
		networkPolicies = append(networkPolicies, &networking.NetworkPolicy{
			ObjectMeta: meta.ObjectMeta{
				Name:      "allow-ingress",
				Namespace: p.Namespace,
			},
			Spec: networking.NetworkPolicySpec{
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"i.4ts.fr/pod": pod.Name,
					},
				},
				PolicyTypes: []networking.PolicyType{
					networking.PolicyTypeIngress,
				},
				Ingress: []networking.NetworkPolicyIngressRule{
					{
						From: []networking.NetworkPolicyPeer{
							{
								NamespaceSelector: &meta.LabelSelector{
									MatchLabels: map[string]string{
										"kubernetes.io/metadata.name": "emissary",
									},
								},
								PodSelector: &meta.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/name": "emissary-ingress-agent",
									},
								},
							},
						},
					},
				},
			},
		})
	}

	networkPolicies = append(networkPolicies, &networking.NetworkPolicy{
		ObjectMeta: meta.ObjectMeta{
			Name:      "allow-egress",
			Namespace: p.Namespace,
		},
		Spec: networking.NetworkPolicySpec{
			PodSelector: meta.LabelSelector{
				MatchLabels: map[string]string{
					"i.4ts.fr/egress": "true",
				},
			},
			PolicyTypes: []networking.PolicyType{
				networking.PolicyTypeEgress,
			},
			Egress: []networking.NetworkPolicyEgressRule{
				{
					To: []networking.NetworkPolicyPeer{
						{
							IPBlock: &networking.IPBlock{
								CIDR: "0.0.0.0/0",
								Except: []string{
									"10.0.0.0/8",
									"192.168.0.0/16",
									"172.16.0.0/20",
								},
							},
						},
					},
				},
			},
		},
	})

	return networkPolicies
}
