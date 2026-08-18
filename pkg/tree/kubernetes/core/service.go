package core

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Service is a core/v1 Service. Only type: LoadBalancer carries cost — it makes
// the cloud-controller-manager provision a real cloud load balancer (ALB/NLB,
// GCP forwarding rule, Azure LB) — so the parser only surfaces those; the Type
// is retained verbatim so downstream consumers can confirm it. ClusterIP /
// NodePort / ExternalName Services are free and are not represented here.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the Service's name and namespace on the embedded
// meta.ObjectMeta; and the Service's Kubernetes labels are stored as the base
// resource's Tags.
type Service struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// Type is spec.type — expected to be "LoadBalancer" for every Service the
	// parser surfaces.
	Type value.String `tree:"type"`

	// Annotations are the Service's Kubernetes annotations, surfaced verbatim.
	// For load balancers these select the flavour that drives the price — e.g.
	// service.beta.kubernetes.io/aws-load-balancer-type (nlb vs the default
	// classic ELB) — as well as the cloud-provider signals the workload type
	// also carries.
	Annotations []resource.Tag `tree:"annotations"`

	// Ports are the service ports (spec.ports), each becoming a listener on the
	// provisioned load balancer.
	Ports []ServicePort `tree:"ports"`
}

// ServicePort is a single entry of a Service's spec.ports.
type ServicePort struct {
	Port     value.Int    `tree:"port"`
	Protocol value.String `tree:"protocol"`
}