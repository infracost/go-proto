package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Subnet struct {
	resource.Resource `tree:"-"`
	AvailabilityZone  value.String `tree:"availability_zone"`

	Relationships SubnetRelationships `tree:"-"`
}

type SubnetRelationships struct {
	NATGateways []*NATGateway
}
