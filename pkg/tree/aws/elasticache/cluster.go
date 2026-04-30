package elasticache

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource          `tree:"-"`
	NodeType                   value.String `tree:"node_type"`
	Engine                     value.Value[CacheEngine] `tree:"engine"`
	EngineVersion              value.String `tree:"engine_version"`
	ParameterGroupName         value.String `tree:"parameter_group_name"`
	ReplicationGroupID         value.String `tree:"replication_group_id"`
	CacheNodes                 value.Int    `tree:"num_cache_nodes"`
	SnapshotRetentionLimitDays value.Int    `tree:"snapshot_retention_limit"`

	Relationships ClusterRelationships `tree:"-"`
}

type ClusterRelationships struct {
	ParameterGroup   *ParameterGroup
	ReplicationGroup *ReplicationGroup
}
