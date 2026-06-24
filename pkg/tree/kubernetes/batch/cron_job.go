package batch

import "github.com/infracost/go-proto/pkg/tree/value"

// CronJob is a batch/v1 CronJob. It wraps a Job template (hence the embedded
// Job) and runs it on a schedule.
type CronJob struct {
	Job      `tree:"-"`
	Schedule value.String `tree:"schedule"`
}
