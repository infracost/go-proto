package gamelift

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type GameServerGroup struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
	MinSize           value.Int    `tree:"min_size"`
	BalancingStrategy value.String `tree:"balancing_strategy"`
}
