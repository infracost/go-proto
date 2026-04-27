package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CapacityProvider struct {
	resource.Resource `tree:"-"`
	Name              value.String                    `tree:"name"`
	ProviderType      value.Value[CapacityProviderType] `tree:"provider_type"`
}

type CapacityProviderType uint32

const (
	CapacityProviderTypeASG              CapacityProviderType = 0
	CapacityProviderTypeManagedInstances CapacityProviderType = 1
)
