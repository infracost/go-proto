// Package batch models the Kubernetes "batch" API group (batch/v1) as a service
// in the tree. Its fields are the workload kinds in that group; each slice is
// tagged with the kind, which becomes the resource Type on the wire.
package batch

// Batch is the batch/v1 API group.
type Batch struct {
	Jobs     []Job     `tree:"job"`
	CronJobs []CronJob `tree:"cronjob"`
}
