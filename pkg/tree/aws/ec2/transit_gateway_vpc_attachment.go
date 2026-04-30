package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TransitGatewayVPCAttachment struct {
	resource.Resource `tree:"-"`
	VPCID             value.String `tree:"vpc_id"`
	TransitGatewayID  value.String `tree:"transit_gateway_id"`

	Relationships TransitGatewayVPCAttachmentRelationships `tree:"-"`
}

type TransitGatewayVPCAttachmentRelationships struct {
	VPC            *VPC
	TransitGateway *TransitGateway
}
