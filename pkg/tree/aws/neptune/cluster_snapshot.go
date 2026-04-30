package neptune

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ClusterSnapshot struct {
	resource.Resource   `tree:"-"`
	DBClusterIdentifier value.String `tree:"db_cluster_identifier"`

	Relationships ClusterSnapshotRelationships `tree:"-"`
}

type ClusterSnapshotRelationships struct {
	Cluster *Cluster
}
