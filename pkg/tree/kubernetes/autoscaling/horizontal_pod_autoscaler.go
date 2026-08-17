package autoscaling

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
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
// against. The metrics driving the autoscaler are not modelled — they decide
// where between the bounds it sits at any moment, which is an observation the
// metrics pipeline reports directly and better.
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

	// Annotations are the HPA's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}
