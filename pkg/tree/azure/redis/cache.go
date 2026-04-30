package redis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cache struct {
	resource.Resource  `tree:"-"`
	SKUName            value.String `tree:"sku_name"`
	Family             value.Value[CacheFamily] `tree:"family"`
	Capacity           value.String `tree:"capacity"`
	ShardCount         value.Int    `tree:"shard_count"`
	ReplicasPerPrimary value.Int    `tree:"replicas_per_primary"`
	ReplicasPerMaster  value.Int    `tree:"replicas_per_master"`
}
