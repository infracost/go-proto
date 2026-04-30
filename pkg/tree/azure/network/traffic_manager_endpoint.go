package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TrafficManagerEndpoint struct {
	resource.Resource `tree:"-"`
	External          value.Bool   `tree:"external"`
	ProfileID         value.String `tree:"profile_id"`

	Relationships TrafficManagerEndpointRelationships `tree:"-"`
}

type TrafficManagerEndpointRelationships struct {
	Profile *TrafficManagerProfile
}
