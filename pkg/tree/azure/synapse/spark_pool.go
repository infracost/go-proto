package synapse

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SparkPool struct {
	resource.Resource `tree:"-"`
	NodeSize          value.Value[SparkPoolNodeSize] `tree:"node_size"`
	NodeCount         value.Int    `tree:"node_count"`
	MinNodeCount      value.Int    `tree:"min_node_count"`
	WorkspaceID       value.String `tree:"workspace_id"`
	NodeCores         value.Int    `tree:"node_cores"`

	Relationships SparkPoolRelationships `tree:"-"`
}

type SparkPoolRelationships struct {
	Workspace *Workspace
}
