// Package kubernetes is the Kubernetes provider in the resource tree. It groups
// workloads by their Kubernetes API group (the service layer): apps/v1 and
// batch/v1. Workloads themselves share a single type — see the workload package.
package kubernetes

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/apps"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/batch"
)

// Kubernetes is the provider node for Kubernetes workloads.
type Kubernetes struct {
	Apps  apps.Apps   `tree:"apps"`
	Batch batch.Batch `tree:"batch"`
}
