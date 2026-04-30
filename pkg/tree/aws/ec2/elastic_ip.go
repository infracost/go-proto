package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ElasticIP struct {
	resource.Resource  `tree:"-"`
	NetworkInterfaceID value.String `tree:"network_interface_id"`
	InstanceID         value.String `tree:"instance_id"`

	Relationships ElasticIPRelationships `tree:"-"`
}

type ElasticIPRelationships struct {
	Association  *ElasticIPAssociation
	NATGateway   *NATGateway
	LoadBalancer *LoadBalancer
}
