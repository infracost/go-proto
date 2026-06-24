package apps

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ReplicaSet is an apps/v1 ReplicaSet. Its replica count is set explicitly via
// spec.replicas.
type ReplicaSet struct {
	workload.Workload `tree:"-"`
	Replicas          value.Int `tree:"replicas"`
}
