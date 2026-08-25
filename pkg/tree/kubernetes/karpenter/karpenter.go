// Package karpenter models the Karpenter CRDs as a service in the tree.
//
// Karpenter provisions nodes directly from pending pods, without a node group
// or an autoscaling group in between. That makes it the one node-provisioning
// path whose definition lives in Kubernetes manifests rather than in Terraform:
// for an EKS managed node group the instance types, sizes and disks are
// Terraform, but under Karpenter they are a NodePool and an EC2NodeClass in the
// cluster's own repository.
//
// It matters twice for rightsizing. First, Karpenter is the regime where a
// per-pod saving actually banks — it consolidates continuously, so freeing
// capacity really does remove nodes, where a fixed-size node group would keep
// running the same instances with more headroom. Second, it is the regime where
// the node-side change can be proposed at all, because these files are the node
// definition and they are locatable in code.
//
// Two API groups sit in this one package. NodePool is karpenter.sh and is
// cloud-agnostic; EC2NodeClass is karpenter.k8s.aws and is the AWS-specific half
// a NodePool points at. They are kept together because neither is meaningful
// alone — a NodePool without its node class does not describe a node — and
// because the other clouds' equivalents will slot in beside them rather than
// forming groups of their own.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps, batch and core groups.
package karpenter

import (
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Karpenter is the group holding the Karpenter CRDs.
type Karpenter struct {
	NodePools      []NodePool     `tree:"nodepool"`
	EC2NodeClasses []EC2NodeClass `tree:"ec2nodeclass"`
}

// Requirement is one entry of a NodePool's spec.template.spec.requirements —
// a constraint on what Karpenter may launch, expressed the same way a node
// selector term is.
//
// This is where instance shape is decided, and it is a constraint rather than a
// choice: a requirement narrows the pool of instance types Karpenter picks from,
// and what it actually launches depends on the pods pending at the time. So
// these bound the cost without determining it, and the instance types a cluster
// really ran are an observation from billing rather than something this file
// states.
type Requirement struct {
	// Key is the label the requirement constrains, e.g.
	// "karpenter.sh/capacity-type" or "node.kubernetes.io/instance-type".
	Key value.String `tree:"key"`

	// Operator is In, NotIn, Exists, DoesNotExist, Gt or Lt.
	Operator value.String `tree:"operator"`

	// Values are the values the operator applies to. Empty for Exists and
	// DoesNotExist, which take none.
	Values *value.List[string] `tree:"values"`
}

// Taint is one entry of a NodePool's spec.template.spec.taints or
// startupTaints — a repulsion that keeps pods off the nodes it launches unless
// they tolerate it.
//
// Relevant to rightsizing because it decides which workloads can land on a pool
// at all: a saving computed against a pool's rate is only real if the workload
// could actually run there.
type Taint struct {
	Key    value.String `tree:"key"`
	Value  value.String `tree:"value"`
	Effect value.String `tree:"effect"`
}
