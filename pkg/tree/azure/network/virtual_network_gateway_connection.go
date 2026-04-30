package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VirtualNetworkGatewayConnection struct {
	resource.Resource          `tree:"-"`
	Type                       value.Value[GatewayConnectionType] `tree:"type"`
	VirtualNetworkGatewayID    value.String                       `tree:"virtual_network_gateway_id"`

	Relationships VirtualNetworkGatewayConnectionRelationships `tree:"-"`
}

type VirtualNetworkGatewayConnectionRelationships struct {
	Gateway *VirtualNetworkGateway
}
