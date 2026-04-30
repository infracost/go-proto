package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TransitGatewayPeeringAttachment struct {
	resource.Resource `tree:"-"`
	TransitGatewayID  value.String `tree:"transit_gateway_id"`

	Relationships TransitGatewayPeeringAttachmentRelationships `tree:"-"`
}

type TransitGatewayPeeringAttachmentRelationships struct {
	TransitGateway *TransitGateway
}
