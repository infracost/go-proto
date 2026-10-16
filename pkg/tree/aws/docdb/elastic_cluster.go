package docdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ElasticCluster struct {
	resource.Resource  `tree:"-"`
	ShardCapacity      value.Int `tree:"shard_capacity"`
	ShardCount         value.Int `tree:"shard_count"`
	ShardInstanceCount value.Int `tree:"shard_instance_count"`
}
