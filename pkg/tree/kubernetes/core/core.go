// Package core models the Kubernetes core API group (the "" group, apiVersion
// v1) as a service in the tree. Unlike the apps and batch groups — whose members
// are all workloads sharing the workload type — the core group holds the
// non-workload kinds.
//
// Two of them are held because they cost money: PersistentVolumeClaim (a
// request for a dynamically provisioned cloud disk) and Service (only type:
// LoadBalancer, which provisions a cloud load balancer). Both are cloud spend
// created outside the cluster's node pools, so nothing in the node-pool IaC
// accounts for them.
//
// The third, Namespace, provisions nothing. It is held so that label/tag
// enforcement can be applied to it — namespaces are where teams record
// ownership, so an unlabelled one is a policy gap that can only be reported if
// the namespace reaches the tree.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps and batch groups.
package core

// Core is the core/v1 API group.
type Core struct {
	PersistentVolumeClaims []PersistentVolumeClaim `tree:"persistentvolumeclaim"`
	Services               []Service               `tree:"service"`
	Namespaces             []Namespace             `tree:"namespace"`
}
