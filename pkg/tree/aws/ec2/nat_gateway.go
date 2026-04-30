package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type NATGateway struct {
	resource.Resource `tree:"-"`
	AllocationID      value.String `tree:"allocation_id"`
	SubnetID          value.String `tree:"subnet_id"`

	Relationships NATGatewayRelationships `tree:"-"`
}

type NATGatewayRelationships struct {
	AllocatedElasticIP *ElasticIP
	Subnet             *Subnet
}
