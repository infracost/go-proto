package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource                `tree:"-"`
	Name                             value.String              `tree:"name"`
	DefaultCapacityProviderStrategy  *CapacityProviderStrategy `tree:"default_capacity_provider_strategy"`

	Relationships ClusterRelationships `tree:"-"`
}

type ClusterRelationships struct {
	CapacityProviders          []*ClusterCapacityProviders
	KnownCapacityProviderTypes map[string]CapacityProviderType
}
