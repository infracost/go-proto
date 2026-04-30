package hdinsight

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type KafkaCluster struct {
	resource.Resource    `tree:"-"`
	HeadNodeVM           value.String `tree:"head_node_vm"`
	WorkerNodeVM         value.String `tree:"worker_node_vm"`
	WorkerInstances      value.Int    `tree:"worker_instances"`
	ZookeeperNodeVM      value.String `tree:"zookeeper_node_vm"`
	Tier                 value.String `tree:"tier"`
	NumberOfDisksPerNode value.Int    `tree:"number_of_disks_per_node"`
}
