package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VPCEndpoint struct {
	resource.Resource `tree:"-"`
	VPCID             value.String                 `tree:"vpc_id"`
	SubnetIDs         value.List[string]           `tree:"subnet_ids"`
	Type              value.Value[VPCEndpointType] `tree:"type"`
	ServiceName       value.String                 `tree:"service_name"`

	Relationships VPCEndpointRelationships `tree:"-"`
}

type VPCEndpointRelationships struct {
	VPC     *VPC
	Subnets []*Subnet
}

type VPCEndpointType uint32

const (
	VPCEndpointTypeUnknown VPCEndpointType = iota
	VPCEndpointTypeGateway
	VPCEndpointTypeInterface
	VPCEndpointTypeGatewayLoadBalancer
)
