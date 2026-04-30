package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VPNGateway struct {
	resource.Resource `tree:"-"`
	ScaleUnits        value.Int                  `tree:"scale_units"`
	Type              value.Value[VPNGatewayType] `tree:"type"`
}
