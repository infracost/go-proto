package karpenter

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ConsolidationPolicy values for a NodePool's
// spec.disruption.consolidationPolicy. This is the setting that decides whether
// a pod-level rightsizing change ever becomes money.
const (
	// ConsolidationPolicyWhenEmptyOrUnderutilized lets Karpenter remove a node
	// whose pods would fit elsewhere, not just an entirely empty one. This is
	// the regime where shrinking requests actually reduces the node count, and
	// it is the default.
	ConsolidationPolicyWhenEmptyOrUnderutilized = "WhenEmptyOrUnderutilized"

	// ConsolidationPolicyWhenEmpty removes a node only once it holds no
	// workload pods at all. Freeing capacity on a partly-used node changes
	// nothing under this policy, so a per-pod saving does not bank until the
	// node happens to empty completely.
	ConsolidationPolicyWhenEmpty = "WhenEmpty"
)

// NodePool is a karpenter.sh NodePool — the constraints Karpenter provisions
// nodes within, and the disruption rules by which it removes them again.
//
// This is the node half of rightsizing, and the half that decides whether the
// pod half is worth anything. Shrinking a container's requests frees capacity;
// capacity becomes money only when a node is removed; and whether a node is
// removed is Disruption below. Under WhenEmptyOrUnderutilized the estate really
// does shrink to fit. Under WhenEmpty it does not, and the same recommendation
// on the same workload saves nothing until a node empties by chance.
//
// The pool does not state what it costs. Requirements bound which instance
// types Karpenter may choose, and it chooses per batch of pending pods, so the
// realised mix is a billing observation rather than something this file
// declares — see Requirement.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the pool's own name and namespace on the embedded
// meta.ObjectMeta; and its Kubernetes labels are stored as the base resource's
// Tags. A NodePool is in fact cluster-scoped, so it has no namespace of its own
// — see the Namespace note on core.Namespace for the same situation.
type NodePool struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// Requirements is spec.template.spec.requirements — the constraints on what
	// Karpenter may launch. An empty list is legal and means unconstrained,
	// which is the widest and most expensive possible pool rather than a
	// misconfiguration.
	Requirements []Requirement `tree:"requirements"`

	// Taints and StartupTaints are spec.template.spec.taints and
	// .startupTaints. Taints persist and gate which workloads may schedule
	// here; startup taints are removed once the node is ready and gate only the
	// boot window.
	Taints        []Taint `tree:"taints"`
	StartupTaints []Taint `tree:"startup_taints"`

	// NodeClassRef is spec.template.spec.nodeClassRef — the EC2NodeClass (or
	// another cloud's equivalent) holding the AMI, disks and networking. A
	// NodePool without it does not describe a node.
	NodeClassRef NodeClassRef `tree:"node_class_ref"`

	// Limits is spec.limits — the ceiling on total resources across every node
	// this pool has running, which is the one hard bound on what it can cost.
	// Unset means unlimited.
	Limits NodePoolLimits `tree:"limits"`

	// Disruption is spec.disruption.
	Disruption Disruption `tree:"disruption"`

	// Annotations are the NodePool's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}

// NodeClassRef is a NodePool's spec.template.spec.nodeClassRef — the pointer to
// the node class that completes it.
type NodeClassRef struct {
	// Group is the API group, e.g. "karpenter.k8s.aws". With Kind it says which
	// cloud's node class this is, and therefore which cloud the pool runs on.
	Group value.String `tree:"group"`

	// Kind is the node class kind, e.g. "EC2NodeClass".
	Kind value.String `tree:"kind"`

	// Name is the node class's metadata.name. Node classes are cluster-scoped,
	// so this alone resolves the reference.
	Name value.String `tree:"name"`
}

// NodePoolLimits is a NodePool's spec.limits — the total resources it may have
// provisioned at once, across all its nodes.
//
// CPU is in millicores and memory in bytes, matching the units the workload
// containers use, so a pool's ceiling and the requests filling it are directly
// comparable. Either may be unset, meaning no limit on that dimension.
type NodePoolLimits struct {
	CPUMillicores value.Int `tree:"cpu_millicores"`
	MemoryBytes   value.Int `tree:"memory_bytes"`
}

// Disruption is a NodePool's spec.disruption — when Karpenter may remove or
// replace the nodes it launched.
type Disruption struct {
	// ConsolidationPolicy is one of the ConsolidationPolicy constants above.
	// Empty means the default, WhenEmptyOrUnderutilized — so an absent policy
	// is the consolidating one, and reading empty as "no consolidation" gets
	// the default backwards.
	ConsolidationPolicy value.String `tree:"consolidation_policy"`

	// ConsolidateAfter is how long a node must sit empty or underutilized
	// before Karpenter acts, as a duration string ("30s", "1h") or "Never".
	// Kept verbatim rather than parsed to seconds, because "Never" is a
	// meaningful value that no number represents.
	ConsolidateAfter value.String `tree:"consolidate_after"`

	// ExpireAfter is spec.template.spec.expireAfter — the node's maximum
	// lifetime before Karpenter replaces it regardless of utilization. Same
	// duration-or-"Never" encoding as ConsolidateAfter.
	ExpireAfter value.String `tree:"expire_after"`

	// Budgets are spec.disruption.budgets — caps on how much disruption may
	// happen at once. A budget of zero over a window blocks consolidation
	// entirely for that period, so like a PodDisruptionBudget it can stop a
	// saving being realised.
	Budgets []DisruptionBudget `tree:"budgets"`
}

// DisruptionBudget is one entry of a Disruption's budgets.
type DisruptionBudget struct {
	// Nodes is how many nodes may be disrupted concurrently, as a count ("3")
	// or a percentage of the pool ("20%"). Kept as a string for the same reason
	// a PodDisruptionBudget's minAvailable is: the two forms mean different
	// things and resolving the percentage needs the live node count.
	Nodes value.String `tree:"nodes"`

	// Schedule and Duration restrict the budget to a window — a cron expression
	// and a duration. Both empty means the budget applies at all times.
	Schedule value.String `tree:"schedule"`
	Duration value.String `tree:"duration"`

	// Reasons limits the budget to particular disruption reasons
	// (Underutilized, Empty, Drifted). Empty means it applies to all of them.
	Reasons *value.List[string] `tree:"reasons"`
}
