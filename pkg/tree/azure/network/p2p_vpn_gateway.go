package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type P2PVPNGateway struct {
	resource.Resource `tree:"-"`
	ScaleUnits        value.Int `tree:"scale_units"`
}
