package hdinsight

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type HBaseCluster struct {
	resource.Resource `tree:"-"`
	HeadNodeVM        value.String `tree:"head_node_vm"`
	RegionNodeVM      value.String `tree:"region_node_vm"`
	RegionInstances   value.Int    `tree:"region_instances"`
	ZookeeperNodeVM   value.String `tree:"zookeeper_node_vm"`
}
