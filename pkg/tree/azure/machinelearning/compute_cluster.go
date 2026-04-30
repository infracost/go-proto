package machinelearning

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ComputeCluster struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
	MinNodeCount      value.Int    `tree:"min_node_count"`
}
