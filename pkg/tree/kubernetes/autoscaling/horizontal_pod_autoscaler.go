package autoscaling

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// MetricSourceType values for a HorizontalPodAutoscaler's spec.metrics[].type.
// The distinction that matters is whether the metric is a share of what the
// container requests: only Resource and ContainerResource are, and only those
// couple the replica count to the numbers a rightsizing recommendation edits.
const (
	// MetricSourceTypeResource scales on a resource the pod requests — cpu or
	// memory — summed across the pod's containers.
	MetricSourceTypeResource = "Resource"

	// MetricSourceTypeContainerResource is the same, narrowed to one named
	// container rather than the pod total. On a multi-container pod this says
	// which container's request the replica count actually keys off.
	MetricSourceTypeContainerResource = "ContainerResource"

	// MetricSourceTypePods scales on a custom per-pod metric averaged over the
	// pods. Not a share of anything requested, so it does not couple to the
	// container's resources.
	MetricSourceTypePods = "Pods"

	// MetricSourceTypeObject scales on a metric describing some other
	// Kubernetes object. The described object is not modelled.
	MetricSourceTypeObject = "Object"

	// MetricSourceTypeExternal scales on a metric from outside the cluster — a
	// queue depth, a request rate. The replica count is then driven by
	// something no manifest describes.
	MetricSourceTypeExternal = "External"
)

// MetricTargetType values for a metric's target.type — which of the target's
// value fields the manifest set.
const (
	// MetricTargetTypeUtilization targets a percentage of the requested
	// resource. This is the setpoint that makes observed headroom expected
	// rather than wasted, and it is only valid on Resource and
	// ContainerResource metrics.
	MetricTargetTypeUtilization = "Utilization"

	// MetricTargetTypeValue targets a raw metric value.
	MetricTargetTypeValue = "Value"

	// MetricTargetTypeAverageValue targets a raw metric value averaged over the
	// pods.
	MetricTargetTypeAverageValue = "AverageValue"
)

// HorizontalPodAutoscaler is an autoscaling/v2 HorizontalPodAutoscaler.
//
// It provisions nothing and costs nothing, and unlike a VerticalPodAutoscaler
// it does not touch container requests — so it does not invalidate a rightsizing
// recommendation. What it invalidates is the workload's declared replica count:
// once an HPA governs a Deployment, spec.replicas in the manifest is read at
// creation and then never again, and a recommendation that proposes editing it
// is proposing a change with no effect.
//
// It also changes how a per-pod saving becomes a real one. Shrinking a request
// on an HPA-governed workload does not reduce the pod count, it makes each pod
// cheaper to schedule and lets the same node fit more of them — so the saving
// only banks if the node count follows, which is the node-coupling question
// rather than the pod one.
//
// MinReplicas and MaxReplicas are held for that reason: they bound how much of
// the estate the workload can occupy, which is what a saving is computed
// against.
//
// Metrics are held for a sharper one. A utilization target is a setpoint rather
// than an observation: an HPA holding a Deployment at 50% CPU produces a
// workload sitting at 50% of its request by design, and a rightsizing pass that
// reads that as half wasted will propose halving the request. Halving it puts
// utilization back at the target, the HPA scales out, and the same spend
// returns as more smaller pods. The observed utilization is something the
// metrics pipeline reports and reports better; the target it is being held at
// exists only in the manifest, the same way a VerticalPodAutoscaler's
// updateMode does.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the HPA's own name and namespace on the embedded
// meta.ObjectMeta; and its Kubernetes labels are stored as the base resource's
// Tags.
type HorizontalPodAutoscaler struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// ScaleTargetRef is spec.scaleTargetRef — the workload this HPA scales.
	// Required by the API, so an empty value means a malformed manifest.
	ScaleTargetRef TargetRef `tree:"scale_target_ref"`

	// MinReplicas is spec.minReplicas. Optional, defaulting to 1 when omitted —
	// so unset is not zero, and reading it as zero would suggest the workload
	// can scale to nothing, which it cannot without a separate feature gate.
	MinReplicas value.Int `tree:"min_replicas"`

	// MaxReplicas is spec.maxReplicas, required by the API. This is the ceiling
	// a worst-case cost is computed against.
	MaxReplicas value.Int `tree:"max_replicas"`

	// Metrics are spec.metrics — what the controller scales on, and the value
	// it holds that signal at. Empty when the manifest states none, in which
	// case the controller falls back to a default CPU utilization target that
	// is cluster configuration rather than repository state.
	Metrics []Metric `tree:"metrics"`

	// Annotations are the HPA's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}

// Metric is one entry of a HorizontalPodAutoscaler's spec.metrics.
//
// The API models this as a five-way union — one nested block per source type,
// each with a target of its own. It is flattened here the way an Ingress path
// flattens its backend: Type says which block the manifest wrote, and the
// fields below carry whichever parts of it mean anything. The described object
// on an Object metric is not modelled; such a metric is recorded so a reader
// knows the replica count is driven from somewhere outside the workload, not so
// that it can be resolved.
type Metric struct {
	// Type is the metric source: one of the MetricSourceType constants above.
	Type value.String `tree:"type"`

	// Name is what the metric is called, which is a different thing per Type.
	// On Resource and ContainerResource it is the resource name — "cpu" or
	// "memory", matching the keys a container's own requests use. On Pods,
	// Object and External it is the custom metric's name, which is arbitrary.
	Name value.String `tree:"name"`

	// ContainerName is the container a ContainerResource metric measures, and
	// empty on every other type. This is the container whose request the
	// replica count keys off, which on a multi-container pod need not be the
	// one a rightsizing recommendation would otherwise pick.
	ContainerName value.String `tree:"container_name"`

	// TargetType is which kind of target the metric states: one of the
	// MetricTargetType constants above.
	TargetType value.String `tree:"target_type"`

	// TargetUtilization is target.averageUtilization as a percentage, so 50
	// means 50%. Set only where TargetType is Utilization, which is the case
	// that decides whether observed headroom is waste or the configuration
	// working as intended.
	//
	// Unset serializes as a zero, and zero is not a target anything states, so
	// read TargetType rather than testing this for absence.
	TargetUtilization value.Int `tree:"target_utilization"`

	// TargetValue is target.value or target.averageValue, whichever TargetType
	// names, kept as the quantity string the manifest wrote. Custom and
	// external metrics carry arbitrary units that nothing here can normalise,
	// so this is not converted to a number the way CPU and memory are.
	TargetValue value.String `tree:"target_value"`
}
