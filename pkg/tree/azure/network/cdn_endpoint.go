package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CDNEndpoint struct {
	resource.Resource    `tree:"-"`
	ProfileName          value.String                    `tree:"profile_name"`
	OptimizationType     value.Value[CDNOptimizationType] `tree:"optimization_type"`
	GlobalDeliveryRules  value.Int                        `tree:"global_delivery_rules"`
	DeliveryRules        value.Int                        `tree:"delivery_rules"`

	Relationships CDNEndpointRelationships `tree:"-"`
}

type CDNEndpointRelationships struct {
	Profile *CDNProfile
}
