package neptune

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Graph is a Neptune Analytics graph.
type Graph struct {
	resource.Resource `tree:"-"`
	ProvisionedMemory value.Int `tree:"provisioned_memory"`
	ReplicaCount      value.Int `tree:"replica_count"`
}
