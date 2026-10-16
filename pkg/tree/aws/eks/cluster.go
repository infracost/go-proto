package eks

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource  `tree:"-"`
	Version            value.String `tree:"version"`
	ControlPlaneTier   value.String `tree:"control_plane_tier"`
	SupportType        value.String `tree:"support_type"`
	HybridNodesEnabled value.Bool   `tree:"hybrid_nodes_enabled"`
}
