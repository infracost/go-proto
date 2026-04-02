package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type GatewayAssociation struct {
	resource.Resource   `tree:"-"`
	AssociatedGatewayID value.String `tree:"associated_gateway_id"`
}
