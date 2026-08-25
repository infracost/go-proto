package policy

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// PodDisruptionBudget is a policy/v1 PodDisruptionBudget.
//
// It provisions nothing and costs nothing, and it is the reason a correct
// rightsizing recommendation can still produce no saving. Freeing capacity on a
// node is only money once the node is removed, removal means draining it, and
// draining evicts pods — which a PDB is entitled to refuse. A budget of
// minAvailable equal to the workload's replica count blocks every voluntary
// eviction outright, and consolidation stalls on it indefinitely.
//
// Unlike the autoscalers, a PDB does not name the workload it protects. It
// selects pods by label, so the join is a selector match rather than a
// reference: a PDB with selector app=web covers every pod carrying that label,
// across however many workloads set it. That makes the relationship many-to-many
// and not resolvable by name, which is why Selector is carried in full rather
// than reduced to a target.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the PDB's own name and namespace on the embedded
// meta.ObjectMeta; and its Kubernetes labels are stored as the base resource's
// Tags. Note that those labels are the PDB's own — Selector below is a
// different thing, the labels it matches pods against.
type PodDisruptionBudget struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// MinAvailable is spec.minAvailable and MaxUnavailable is
	// spec.maxUnavailable. Exactly one of the two is set on a valid PDB; both
	// empty means a malformed manifest.
	//
	// Both are Kubernetes IntOrString: "2" and "50%" are equally valid and mean
	// different things, the percentage being resolved against the number of
	// pods the selector currently matches. They are kept as strings so the
	// distinction survives — parsing "50%" to 50 would read as an absolute
	// count and silently invert the meaning on any workload with more than a
	// hundred pods.
	MinAvailable   value.String `tree:"min_available"`
	MaxUnavailable value.String `tree:"max_unavailable"`

	// Selector is spec.selector.matchLabels — the labels a pod must carry for
	// this budget to cover it. Held as tags so it reuses the same machinery as
	// the label sets elsewhere in the tree, and so a consumer can compare it
	// against a workload's labels directly.
	//
	// Only matchLabels is modelled. spec.selector.matchExpressions is the more
	// general form (In, NotIn, Exists, DoesNotExist over a key) and cannot be
	// expressed as key/value pairs; a PDB using it will have an empty or partial
	// Selector here, so an empty Selector means "not expressible" rather than
	// "matches nothing".
	Selector []resource.Tag `tree:"selector"`

	// Annotations are the PDB's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}
