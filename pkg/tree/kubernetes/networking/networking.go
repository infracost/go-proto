// Package networking models the Kubernetes "networking.k8s.io" API group
// (networking.k8s.io/v1) as a service in the tree.
//
// Its member, Ingress, is here for the same reason core holds a LoadBalancer
// Service: it provisions billable cloud infrastructure that nothing in the node
// pool IaC accounts for. An Ingress handled by the AWS Load Balancer Controller
// creates a real ALB, priced per hour plus per LCU, and the manifest that
// creates it is the only place that spend is declared.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps, batch and core groups.
package networking

// Networking is the networking.k8s.io/v1 API group.
type Networking struct {
	Ingresses []Ingress `tree:"ingress"`
}
