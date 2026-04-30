package neptune

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ClusterInstance struct {
	resource.Resource `tree:"-"`
	InstanceClass     value.String `tree:"instance_class"`
	ClusterIdentifier value.String `tree:"cluster_identifier"`

	Relationships ClusterInstanceRelationships `tree:"-"`
}

type ClusterInstanceRelationships struct {
	Cluster *Cluster
}
