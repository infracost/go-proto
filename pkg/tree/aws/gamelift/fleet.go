package gamelift

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Fleet struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
	Spot              value.Bool   `tree:"spot"`
	Anywhere          value.Bool   `tree:"anywhere"`
	BuildID           value.String `tree:"build_id"`
	FromScript        value.Bool   `tree:"from_script"`

	Relationships FleetRelationships `tree:"-"`
}

type FleetRelationships struct {
	Build *Build
}
