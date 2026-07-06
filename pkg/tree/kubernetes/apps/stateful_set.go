package apps

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/core"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// StatefulSet is an apps/v1 StatefulSet. Its replica count is set explicitly via
// spec.replicas.
type StatefulSet struct {
	workload.Workload `tree:"-"`
	Replicas          value.Int `tree:"replicas"`

	// VolumeClaimTemplates are spec.volumeClaimTemplates. A StatefulSet
	// provisions one persistent volume per template per replica, so the storage
	// cost is (sum of these) x Replicas — which is why they live on the
	// StatefulSet rather than as standalone PersistentVolumeClaims in the core
	// group.
	VolumeClaimTemplates []core.StorageRequest `tree:"volume_claim_template"`
}
