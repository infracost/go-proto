package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	Location          value.String `tree:"location"`
	AutopilotEnabled  value.Bool   `tree:"autopilot_enabled"`
	Zones             value.Int    `tree:"zones"`
	IsZone            value.Bool   `tree:"is_zone"`

	Relationships ClusterRelationships `tree:"-"`
}

type ClusterRelationships struct {
	DefaultNodePool *NodePool
	NodePools       []*NodePool
}
