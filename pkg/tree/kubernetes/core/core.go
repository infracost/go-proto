// Package core models the Kubernetes core API group (the "" group, apiVersion
// v1) as a service in the tree. Unlike the apps and batch groups — whose members
// are all workloads sharing the workload type — the core group holds the two
// cost-relevant non-workload kinds: PersistentVolumeClaim (a request for a
// dynamically provisioned cloud disk) and Service (only type: LoadBalancer,
// which provisions a cloud load balancer). Both are cloud spend created outside
// the cluster's node pools, so nothing in the node-pool IaC accounts for them.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps and batch groups.
package core

// Core is the core/v1 API group.
type Core struct {
	PersistentVolumeClaims []PersistentVolumeClaim `tree:"persistentvolumeclaim"`
	Services               []Service               `tree:"service"`
}