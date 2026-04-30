package hdinsight

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type HadoopCluster struct {
	resource.Resource `tree:"-"`
	HeadNodeVM        value.String `tree:"head_node_vm"`
	WorkerNodeVM      value.String `tree:"worker_node_vm"`
	ZookeeperNodeVM   value.String `tree:"zookeeper_node_vm"`
	EdgeNodeVM        value.String `tree:"edge_node_vm"`
	WorkerInstances   value.Int    `tree:"worker_instances"`
	EdgeInstances     value.Int    `tree:"edge_instances"`
}
