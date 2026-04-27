package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ClusterCapacityProviders struct {
	resource.Resource                 `tree:"-"`
	Name                              value.String               `tree:"name"`
	ClusterName                       value.String               `tree:"cluster_name"`
	CapacityProviders                 []value.String             `tree:"capacity_providers"`
	DefaultCapacityProviderStrategies []CapacityProviderStrategy `tree:"default_capacity_provider_strategies"`
}

type CapacityProviderStrategy struct {
	CapacityProvider value.String `tree:"capacity_provider"`
	Weight           value.Int    `tree:"weight"`
	Base             value.Int    `tree:"base"`
}
