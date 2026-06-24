package apps

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// StatefulSet is an apps/v1 StatefulSet. Its replica count is set explicitly via
// spec.replicas.
type StatefulSet struct {
	workload.Workload `tree:"-"`
	Replicas          value.Int `tree:"replicas"`
}
