// Package kubernetes is the Kubernetes provider in the resource tree. It groups
// resources by their Kubernetes API group (the service layer): apps/v1,
// batch/v1 and the core (v1) group. The apps and batch groups hold workloads,
// which share a single type — see the workload package. The core group holds
// the cost-relevant non-workload kinds (PersistentVolumeClaim, LoadBalancer
// Service) — see the core package.
package kubernetes

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/apps"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/batch"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/core"
)

// Kubernetes is the provider node for Kubernetes resources.
type Kubernetes struct {
	Apps  apps.Apps   `tree:"apps"`
	Batch batch.Batch `tree:"batch"`
	Core  core.Core   `tree:"core"`
}
