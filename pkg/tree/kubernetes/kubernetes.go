// Package kubernetes is the Kubernetes provider in the resource tree. It groups
// resources by their Kubernetes API group (the service layer): apps/v1,
// batch/v1 and the core (v1) group. The apps and batch groups hold workloads,
// which share a single type — see the workload package. The core group holds
// the cost-relevant non-workload kinds (PersistentVolumeClaim, LoadBalancer
// Service) — see the core package.
//
// Not every group is here because it costs money. The autoscaling and policy
// groups provision nothing at all: they are held because they decide whether a
// change to a workload's manifest takes effect, and whether the capacity it
// frees ever becomes a smaller bill. A tree that models only what is billable
// can price a cluster but cannot tell whether a proposed change to it would
// work, which is what those groups are for — see their package docs.
//
// The karpenter group is the one that spans both. It is a set of CRDs rather
// than a built-in API group, and it is the only place a node definition lives
// in Kubernetes manifests rather than in Terraform.
package kubernetes

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/apps"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/autoscaling"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/batch"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/core"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/karpenter"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/networking"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/policy"
)

// Kubernetes is the provider node for Kubernetes resources.
type Kubernetes struct {
	Apps        apps.Apps               `tree:"apps"`
	Batch       batch.Batch             `tree:"batch"`
	Core        core.Core               `tree:"core"`
	Autoscaling autoscaling.Autoscaling `tree:"autoscaling"`
	Policy      policy.Policy           `tree:"policy"`
	Networking  networking.Networking   `tree:"networking"`
	Karpenter   karpenter.Karpenter     `tree:"karpenter"`
}
