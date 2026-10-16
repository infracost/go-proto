package appstream

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Fleet struct {
	resource.Resource      `tree:"-"`
	InstanceType           value.String           `tree:"instance_type"`
	FleetType              value.Value[FleetType] `tree:"fleet_type"`
	DesiredInstances       value.Int              `tree:"desired_instances"`
	DesiredSessions        value.Int              `tree:"desired_sessions"`
	MaxSessionsPerInstance value.Int              `tree:"max_sessions_per_instance"`
}

type FleetType uint32

const (
	FleetTypeUnknown FleetType = iota
	FleetTypeOnDemand
	FleetTypeAlwaysOn
	FleetTypeElastic
)
