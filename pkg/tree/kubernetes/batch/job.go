package batch

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Job is a batch/v1 Job. Completions and Parallelism govern how many pods run,
// which the downstream cost step uses alongside the per-container sizing.
type Job struct {
	workload.Workload `tree:"-"`
	Completions       value.Int `tree:"completions"`
	Parallelism       value.Int `tree:"parallelism"`
}
