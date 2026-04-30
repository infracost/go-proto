package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type AppEnvironment struct {
	resource.Resource      `tree:"-"`
	InfrastructureSubnetID value.String       `tree:"infrastructure_subnet_id"`
	ZoneRedundant          value.Bool         `tree:"zone_redundant"`
	WorkloadProfiles       []WorkloadProfile `tree:"workload_profiles"`
}

type WorkloadProfile struct {
	Name         value.String `tree:"name"`
	Type         value.String `tree:"type"`
	MinimumCount value.Int    `tree:"minimum_count"`
	MaximumCount value.Int    `tree:"maximum_count"`
}
