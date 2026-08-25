package meta

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// LabelSelector is a Kubernetes meta/v1 LabelSelector — the way one object
// states which pods it applies to, rather than naming them.
//
// The two halves are a conjunction, and so is each half: a pod matches when it
// carries every label in MatchLabels and satisfies every expression in
// MatchExpressions. An empty selector matches everything, which is the opposite
// of matching nothing and is worth stating because the mistake is easy.
//
// Both halves are carried because they are not interchangeable. MatchLabels can
// only say that a key equals a value; MatchExpressions can say that a key is one
// of several values, is none of them, is present, or is absent. Reducing the
// second to the first is only possible for In over a single value, and dropping
// the rest widens the selector — a consumer asking "does this cover my workload"
// then gets a superset, and one asking "which workloads does this cover" gets a
// confident wrong answer.
type LabelSelector struct {
	// MatchLabels is spec.selector.matchLabels — labels a pod must carry
	// exactly. Held as tags so it reuses the same machinery as the label sets
	// elsewhere in the tree, and so a consumer can compare it against a
	// workload's pod labels directly.
	MatchLabels []resource.Tag `tree:"match_labels"`

	// MatchExpressions is spec.selector.matchExpressions — the terms that state
	// something a key/value pair cannot.
	MatchExpressions []LabelSelectorRequirement `tree:"match_expressions"`
}

// LabelSelectorRequirement is one term of a LabelSelector's matchExpressions.
//
// This is deliberately a separate type from karpenter.Requirement, which has the
// same three fields. Karpenter's is a NodeSelectorRequirement, whose operator set
// also includes Gt and Lt; a shared type would imply those are valid on a label
// selector, which they are not.
type LabelSelectorRequirement struct {
	// Key is the label the term constrains.
	Key value.String `tree:"key"`

	// Operator is In, NotIn, Exists or DoesNotExist.
	Operator value.String `tree:"operator"`

	// Values are the values the operator applies to. Nil for Exists and
	// DoesNotExist, which take none — a distinction that has to survive, since
	// an empty list under In matches nothing while Exists matches any value.
	Values *value.List[string] `tree:"values"`
}
