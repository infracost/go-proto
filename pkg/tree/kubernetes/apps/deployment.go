package apps

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Deployment is an apps/v1 Deployment. Its replica count is set explicitly via
// spec.replicas.
type Deployment struct {
	workload.Workload `tree:"-"`
	Replicas          value.Int `tree:"replicas"`
}
