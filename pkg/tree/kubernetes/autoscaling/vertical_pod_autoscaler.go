package autoscaling

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// UpdateMode values for a VerticalPodAutoscaler's spec.updatePolicy.updateMode.
// The distinction that matters is whether the controller writes requests back
// onto pods: only Off leaves the manifest's values in force.
const (
	// UpdateModeOff computes recommendations and does nothing with them. The
	// manifest's requests are what the pods run with, so editing the manifest
	// works normally — a VPA in this mode is an observation, not an obstacle.
	UpdateModeOff = "Off"

	// UpdateModeInitial applies recommendations at pod creation only. Existing
	// pods keep their current requests; the next rollout takes the
	// controller's numbers rather than the manifest's.
	UpdateModeInitial = "Initial"

	// UpdateModeRecreate applies recommendations at pod creation and evicts
	// running pods whose requests drift far enough from them.
	UpdateModeRecreate = "Recreate"

	// UpdateModeAuto is the default when updatePolicy is omitted entirely, and
	// currently behaves as Recreate. Note the consequence for parsing: a VPA
	// with no updatePolicy block is not "unset" but "Auto", so treating an
	// absent mode as harmless gets the common case backwards.
	UpdateModeAuto = "Auto"
)

// ContainerPolicyMode values for a container policy's mode field.
const (
	// ContainerPolicyModeAuto applies the VPA's recommendations to this
	// container. The default when a container policy omits mode.
	ContainerPolicyModeAuto = "Auto"

	// ContainerPolicyModeOff exempts this container, so its manifest requests
	// stand even when the VPA as a whole is in an applying mode. A workload can
	// therefore be partly governed: a sidecar left to the VPA and the app
	// container opted out, or the reverse.
	ContainerPolicyModeOff = "Off"
)

// ControlledValues values for a container policy's controlledValues field — how
// far the controller's writes reach into the container's resources.
const (
	// ControlledValuesRequestsAndLimits lets the controller write both requests
	// and limits, scaling the limit to preserve the ratio the manifest
	// declared. This is the default when the field is omitted, so an absent
	// controlledValues means the manifest's limits are overwritten too.
	ControlledValuesRequestsAndLimits = "RequestsAndLimits"

	// ControlledValuesRequestsOnly leaves limits alone. An edit to a
	// container's limits then survives, though an edit to its requests still
	// does not.
	ControlledValuesRequestsOnly = "RequestsOnly"
)

// VerticalPodAutoscaler is an autoscaling.k8s.io VerticalPodAutoscaler — a CRD
// installed alongside the VPA controller, not a built-in kind.
//
// It provisions nothing and costs nothing. It is in the tree because in every
// mode but Off it sets container requests at admission, which means a
// recommendation that edits the workload's manifest changes a value the cluster
// then overwrites. The pull request merges, the pods keep their old sizes, and
// the finding never resolves. Nothing in the metrics reports that a VPA exists,
// so the manifest is the only place this is visible.
//
// Reading it in the other direction is the useful one: given a workload, is
// there a VPA whose TargetRef names it, and if so does anything exempt the
// container and resource being recommended on — UpdateMode Off for the VPA as a
// whole, or a container policy narrower than that. It decides whether the
// workload is fixable by editing code at all, or whether the recommendation
// belongs in the VPA's own resource policy instead.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the VPA's own name and namespace on the embedded
// meta.ObjectMeta; and its Kubernetes labels are stored as the base resource's
// Tags.
type VerticalPodAutoscaler struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// TargetRef is spec.targetRef — the workload this VPA governs. Required by
	// the CRD, so an empty value means the manifest is malformed rather than
	// that the VPA applies broadly.
	TargetRef TargetRef `tree:"target_ref"`

	// UpdateMode is spec.updatePolicy.updateMode: one of the UpdateMode
	// constants above.
	//
	// An absent updatePolicy means Auto, not "no mode" — the emptiest possible
	// VPA manifest is one of the applying ones. Whether the parser records that
	// default explicitly or leaves the value unset for the consumer to
	// interpret is the parser's decision; either way an empty value must not be
	// read as Off.
	UpdateMode value.String `tree:"update_mode"`

	// ContainerPolicies is spec.resourcePolicy.containerPolicies, empty when the
	// VPA states no per-container policy. Each entry can exempt a container or
	// bound what the controller is allowed to set for it.
	ContainerPolicies []ContainerPolicy `tree:"container_policies"`

	// Annotations are the VPA's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}

// ContainerPolicy is one entry of a VerticalPodAutoscaler's
// spec.resourcePolicy.containerPolicies — the per-container overrides on what
// the controller may do.
//
// These matter to a recommendation three times over. Mode decides whether a
// given container is governed at all. ControlledResources and ControlledValues
// narrow that to which of its numbers are governed, so a container can be
// governed for memory and left alone for CPU. And the Min/Max bounds are the
// range the controller will keep it inside — a recommendation outside those
// bounds cannot take effect even where the VPA is applying.
type ContainerPolicy struct {
	// ContainerName is the container this policy applies to. The wildcard "*"
	// matches every container in the pod, and is how a whole workload is
	// exempted with a single entry — so this is not always a real container
	// name.
	ContainerName value.String `tree:"container_name"`

	// Mode is the policy's mode: one of the ContainerPolicyMode constants
	// above. Empty means Auto, the CRD's default.
	Mode value.String `tree:"mode"`

	// ControlledResources is the policy's controlledResources — which of the
	// container's resources the controller sets, by resource name ("cpu",
	// "memory").
	//
	// This is Mode narrowed to a single resource: a policy listing only
	// "memory" leaves CPU requests alone, so a CPU rightsizing recommendation
	// on that container is actionable while a memory one is not. Nil means the
	// field was omitted, which defaults to both — so absence is the governing
	// case, and must not be read as controlling nothing.
	ControlledResources *value.List[string] `tree:"controlled_resources"`

	// ControlledValues is the policy's controlledValues: one of the
	// ControlledValues constants above. Empty means RequestsAndLimits, the
	// CRD's default — so as with UpdateMode an absent value is the wider of the
	// two, and reading empty as requests-only gets it backwards.
	ControlledValues value.String `tree:"controlled_values"`

	// MinAllowed and MaxAllowed bound what the controller may set. Either side
	// of either pair may be unset, meaning unbounded in that direction — which
	// is why these are not plain numbers; see ResourceAmounts.
	MinAllowed ResourceAmounts `tree:"min_allowed"`
	MaxAllowed ResourceAmounts `tree:"max_allowed"`
}
