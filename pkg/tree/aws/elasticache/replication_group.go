package elasticache

import (
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ReplicationGroup struct {
	resource.Resource               `tree:"-"`
	ID                              value.String `tree:"replication_group_id"`
	NodeType                        value.String `tree:"node_type"`
	Engine                          value.Value[CacheEngine] `tree:"engine"`
	EngineVersion                   value.String `tree:"engine_version"`
	CacheClusterCount               value.Int    `tree:"num_cache_clusters"`
	ClusterNodeGroupCount           value.Int    `tree:"num_node_groups"`
	ClusterReplicaCountPerNodeGroup value.Int    `tree:"replicas_per_node_group"`
	SnapshotRetentionLimitDays      value.Int    `tree:"snapshot_retention_limit"`
	ParameterGroupName              value.String `tree:"parameter_group_name"`

	Relationships ReplicationGroupRelationships `tree:"-"`
}

type ReplicationGroupRelationships struct {
	AppAutoscalingTargets []*appautoscaling.Target
}
