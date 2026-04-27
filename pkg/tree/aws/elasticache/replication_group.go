package elasticache

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ReplicationGroup struct {
	resource.Resource           `tree:"-"`
	ReplicationGroupID          value.String `tree:"replication_group_id"`
	NodeType                    value.String `tree:"node_type"`
	Engine                      value.String `tree:"engine"`
	EngineVersion               value.String `tree:"engine_version"`
	CacheClusters               value.Int    `tree:"num_cache_clusters"`
	ClusterNodeGroups           value.Int    `tree:"num_node_groups"`
	ClusterReplicasPerNodeGroup value.Int    `tree:"replicas_per_node_group"`
	SnapshotRetentionLimit      value.Int    `tree:"snapshot_retention_limit"`
	ParameterGroupName          value.String `tree:"parameter_group_name"`
}
