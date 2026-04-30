package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type GatewayAssociation struct {
	resource.Resource   `tree:"-"`
	AssociatedGatewayID value.String `tree:"associated_gateway_id"`

	Relationships GatewayAssociationRelationships `tree:"-"`
}

type GatewayAssociationRelationships struct {
	TransitGateway *ec2.TransitGateway
}
