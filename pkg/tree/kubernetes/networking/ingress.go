package networking

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Ingress is a networking.k8s.io/v1 Ingress.
//
// Whether it costs anything depends entirely on who handles it, and the
// manifest says that only indirectly. An ingress-nginx Ingress is a routing rule
// on a load balancer that already exists — free, and already counted as the
// controller's own Service. An Ingress handled by the AWS Load Balancer
// Controller provisions a dedicated ALB per Ingress by default, priced per hour
// plus per LCU for as long as it exists.
//
// The controller is chosen by IngressClassName, and its behaviour is configured
// by Annotations — including the one that decides whether several Ingresses
// share one ALB or get one each, which is the difference between a single
// hourly charge and one per Ingress. Both are carried verbatim rather than
// interpreted here: the mapping from class and annotation to a priced resource
// is provider knowledge, and this package stays provider-agnostic the way the
// workload annotations do.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the Ingress's own name and namespace on the
// embedded meta.ObjectMeta; and its Kubernetes labels are stored as the base
// resource's Tags.
type Ingress struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// IngressClassName is spec.ingressClassName — which controller handles this
	// Ingress, and therefore whether it provisions anything.
	//
	// Optional, and an empty value is not the same as "none": with it omitted
	// the cluster's default IngressClass applies, which is cluster state and not
	// in the repository. So an empty class means the handler is unknown rather
	// than that there is no handler, and it cannot be read as free. The older
	// kubernetes.io/ingress.class annotation carries the same information and
	// still appears in the wild; it arrives in Annotations below.
	IngressClassName value.String `tree:"ingress_class_name"`

	// Annotations are the Ingress's Kubernetes annotations, surfaced verbatim.
	// This is where the load-balancer configuration lives for every controller
	// that provisions one — scheme, target type, and the group that decides ALB
	// sharing — so for cost purposes they carry more than the spec does.
	Annotations []resource.Tag `tree:"annotations"`

	// Rules are spec.rules. They do not affect what is provisioned — one ALB
	// serves all of them — but they are the join to the Services behind it, and
	// the host names are how a reader recognises what the load balancer is for.
	Rules []IngressRule `tree:"rules"`
}

// IngressRule is one entry of an Ingress's spec.rules: a host and the paths
// routed under it.
type IngressRule struct {
	// Host is the rule's host. Optional — a rule with no host matches any
	// request not matched by a hosted rule, which is the catch-all form.
	Host value.String `tree:"host"`

	// Paths are the rule's http.paths entries.
	Paths []IngressPath `tree:"paths"`
}

// IngressPath is one entry of an IngressRule's http.paths — a path and the
// Service it routes to.
type IngressPath struct {
	// Path is the matched path, e.g. "/api".
	Path value.String `tree:"path"`

	// PathType is Exact, Prefix or ImplementationSpecific.
	PathType value.String `tree:"path_type"`

	// ServiceName and ServicePort are backend.service.name and
	// backend.service.port — the in-cluster Service this path routes to.
	//
	// A backend can instead be resource (an ObjectRef to a cluster resource),
	// in which case both are empty. The port may be named rather than numeric,
	// so ServicePort is a string.
	ServiceName value.String `tree:"service_name"`
	ServicePort value.String `tree:"service_port"`
}
