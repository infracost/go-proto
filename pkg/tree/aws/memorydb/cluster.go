package memorydb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource   `tree:"-"`
	NodeType            value.String        `tree:"node_type"`
	Engine              value.Value[Engine] `tree:"engine"`
	NumShards           value.Int           `tree:"num_shards"`
	NumReplicasPerShard value.Int           `tree:"num_replicas_per_shard"`
}

type Engine uint32

const (
	EngineUnknown Engine = iota
	EngineRedis
	EngineValkey
)
