// Package policy models the Kubernetes "policy" API group (policy/v1) as a
// service in the tree.
//
// Its one member, PodDisruptionBudget, provisions nothing and costs nothing. It
// is held because it decides whether a saving can be realised: a rightsizing
// change frees capacity on a node, and that capacity only becomes money when
// the node itself goes away. Removing a node means draining it, and a
// PodDisruptionBudget can refuse the eviction that drain depends on — so a
// workload can be correctly rightsized and produce no saving at all.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps, batch and core groups.
package policy

// Policy is the policy/v1 API group.
type Policy struct {
	PodDisruptionBudgets []PodDisruptionBudget `tree:"poddisruptionbudget"`
}
