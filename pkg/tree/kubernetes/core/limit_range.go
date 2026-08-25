package core

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// LimitRangeType values for a LimitRangeItem's type — the object a limit
// applies to.
const (
	// LimitRangeTypeContainer applies per container. The type that supplies
	// default requests and limits, and so the one that decides what a container
	// declaring none actually runs with.
	LimitRangeTypeContainer = "Container"

	// LimitRangeTypePod applies to the sum across a pod's containers. It
	// constrains but supplies no defaults.
	LimitRangeTypePod = "Pod"

	// LimitRangeTypePersistentVolumeClaim bounds the storage a claim in this
	// namespace may request.
	LimitRangeTypePersistentVolumeClaim = "PersistentVolumeClaim"
)

// LimitRange is a core/v1 LimitRange — namespace-wide defaults and bounds on
// what containers, pods and claims may request.
//
// It provisions nothing, and it is here because it makes a manifest an
// incomplete account of what a workload runs with. A container that declares no
// CPU request is not therefore unrequested: if its namespace has a LimitRange
// with a defaultRequest, the admission controller writes that value in, and the
// pod runs with a request no file in the repository states.
//
// That decides whether a workload is a finding at all. The recommendation
// engine treats a nil request as a case it cannot compute a ratio for and has
// to handle separately — but where a LimitRange supplies one, the request is
// known, the ratio is computable, and the workload is an ordinary candidate.
// Without this kind in the tree the two situations are indistinguishable.
//
// The bounds matter in the other direction: Min and Max are what admission will
// accept, so a recommendation below Min is one the cluster rejects, and the pod
// fails to schedule rather than running smaller.
//
// A LimitRange applies by namespace and names no workload, so the join is
// scope-based rather than a reference — every container in the namespace is
// governed by it, and the namespace on the embedded meta.ObjectMeta is the whole
// of the relationship.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the LimitRange's own name and namespace on the
// embedded meta.ObjectMeta; and its Kubernetes labels are stored as the base
// resource's Tags.
type LimitRange struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// Limits is spec.limits. A LimitRange with no entries constrains nothing;
	// several entries of different types commonly appear together.
	Limits []LimitRangeItem `tree:"limits"`

	// Annotations are the LimitRange's Kubernetes annotations, surfaced
	// verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}

// LimitRangeItem is one entry of a LimitRange's spec.limits.
//
// Every field is optional, and the difference between unset and zero is real
// here — a zero default request means "request nothing", where absent means
// "supply nothing and leave the container's own value in force". Absence is not
// representable on its own, though: an unset value serializes as a zero. Where
// the distinction matters the parser marks the value it filled in with
// flag.Synthetic, which does survive, the same way meta.ObjectMeta.Namespace
// records an assumed namespace.
type LimitRangeItem struct {
	// Type is one of the LimitRangeType constants above.
	Type value.String `tree:"type"`

	// DefaultRequest is defaultRequest — what a container gets when it declares
	// no request. This is the field that makes an apparently unrequested
	// container a known quantity, and the one a recommendation has to read
	// before concluding a workload has no request to rightsize.
	DefaultRequest LimitRangeAmounts `tree:"default_request"`

	// Default is default — the limit a container gets when it declares none.
	// Note the asymmetry in the Kubernetes API: the field named "default" sets
	// limits, and requests are set by the separate "defaultRequest" above.
	Default LimitRangeAmounts `tree:"default"`

	// Min and Max are the bounds admission enforces. A pod requesting outside
	// them is rejected rather than clamped, so a recommendation below Min does
	// not make the workload smaller — it stops it scheduling.
	Min LimitRangeAmounts `tree:"min"`
	Max LimitRangeAmounts `tree:"max"`

	// MaxLimitRequestRatio caps how far a limit may exceed its request, per
	// resource. It bounds overcommit, and it can make a recommendation
	// unschedulable on its own: lowering a request without lowering the limit
	// widens the ratio, and admission rejects the pod once it crosses this.
	MaxLimitRequestRatio LimitRangeRatios `tree:"max_limit_request_ratio"`
}

// LimitRangeAmounts is a CPU/memory/storage triple as a LimitRangeItem states
// it, in the base units the rest of the tree uses — CPU in millicores, memory
// and storage in bytes.
//
// Storage appears alongside the other two because a PersistentVolumeClaim-typed
// item bounds storage rather than compute; a Container-typed item leaves it
// unset.
type LimitRangeAmounts struct {
	CPUMillicores value.Int `tree:"cpu_millicores"`
	MemoryBytes   value.Int `tree:"memory_bytes"`
	StorageBytes  value.Int `tree:"storage_bytes"`
}

// LimitRangeRatios is a LimitRangeItem's maxLimitRequestRatio — a bare multiple
// per resource, not a quantity, so it carries no units and is a ratio rather
// than an amount.
type LimitRangeRatios struct {
	CPU    value.Double `tree:"cpu"`
	Memory value.Double `tree:"memory"`
}
