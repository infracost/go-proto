// Package apps models the Kubernetes "apps" API group (apps/v1) as a service in
// the tree. Its fields are the workload kinds in that group; each slice is
// tagged with the kind, which becomes the resource Type on the wire.
package apps

import "github.com/infracost/go-proto/pkg/tree/kubernetes/workload"

// Apps is the apps/v1 API group. DaemonSet has no replica count — it runs one
// pod per node — so it uses the bare workload type.
type Apps struct {
	Deployments  []Deployment        `tree:"deployment"`
	StatefulSets []StatefulSet       `tree:"statefulset"`
	DaemonSets   []workload.Workload `tree:"daemonset"`
	ReplicaSets  []ReplicaSet        `tree:"replicaset"`
}
