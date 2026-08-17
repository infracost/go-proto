// Package autoscaling models the Kubernetes autoscalers as a service in the
// tree. Unlike the apps, batch and core groups, its members provision nothing
// and cost nothing directly — they are held because they decide whether a
// change to a workload's manifest has any effect.
//
// A rightsizing recommendation edits the requests and limits declared on a
// workload. An autoscaler attached to that workload can overwrite those values
// at admission (VerticalPodAutoscaler) or make the declared replica count dead
// config (HorizontalPodAutoscaler). Either way a pull request that edits the
// manifest changes a number the cluster ignores, and the finding never
// resolves. That is a correctness problem rather than a coverage one, which is
// why these are in the tree despite carrying no cost of their own.
//
// The two kinds are separate Kubernetes API groups — HorizontalPodAutoscaler is
// autoscaling/v2, built in; VerticalPodAutoscaler is autoscaling.k8s.io/v1, a
// CRD installed with the VPA controller. They share this package because they
// are the same concern from the tree's point of view: something other than the
// manifest is deciding a workload's resources.
//
// Each slice is tagged with the kind, which becomes the resource Type on the
// wire — mirroring the apps, batch and core groups.
package autoscaling

import (
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Autoscaling is the group holding the Kubernetes autoscaler kinds.
type Autoscaling struct {
	VerticalPodAutoscalers   []VerticalPodAutoscaler   `tree:"verticalpodautoscaler"`
	HorizontalPodAutoscalers []HorizontalPodAutoscaler `tree:"horizontalpodautoscaler"`
}

// TargetRef identifies the workload an autoscaler acts on — spec.targetRef on a
// VerticalPodAutoscaler, spec.scaleTargetRef on a HorizontalPodAutoscaler.
//
// This is the join back to the workload, and it is a pointer rather than an
// identity: the autoscaler names the workload, so finding the autoscaler that
// governs a given Deployment means searching by these fields rather than by the
// autoscaler's own name, which is arbitrary and frequently unrelated.
//
// Kind and Name are namespace-scoped — an autoscaler can only target a workload
// in its own namespace — so the namespace on the embedding kind's ObjectMeta
// completes the reference.
type TargetRef struct {
	// APIVersion is the target's apiVersion, e.g. "apps/v1". Optional in both
	// CRDs and frequently omitted, so an empty value means unspecified rather
	// than absent.
	APIVersion value.String `tree:"api_version"`

	// Kind is the target's kind, e.g. "Deployment" or "StatefulSet". Recorded
	// verbatim from the manifest, so it keeps the CamelCase the Kubernetes API
	// uses rather than the lower-cased form the tree's addresses carry.
	Kind value.String `tree:"kind"`

	// Name is the target workload's metadata.name.
	Name value.String `tree:"name"`
}

// ResourceAmounts is a CPU/memory pair as an autoscaler or a LimitRange states
// it, in the same base units the workload containers use — CPU in millicores,
// memory in bytes — so a bound can be compared against a container's request
// without converting first.
//
// Both are optional in every position they appear: a policy may bound only CPU,
// or only memory. Note that absence is not representable on its own — an unset
// value serializes as a zero — and a zero bound reads as "pinned to nothing"
// rather than "unbounded". Where the distinction matters the parser marks the
// value it filled in with flag.Synthetic, which does survive, the same way
// meta.ObjectMeta.Namespace records an assumed namespace.
type ResourceAmounts struct {
	CPUMillicores value.Int `tree:"cpu_millicores"`
	MemoryBytes   value.Int `tree:"memory_bytes"`
}
