package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VirtualNetworkPeering struct {
	resource.Resource        `tree:"-"`
	VirtualNetworkName       value.String `tree:"virtual_network_name"`
	RemoteVirtualNetworkID   value.String `tree:"remote_virtual_network_id"`

	Relationships VirtualNetworkPeeringRelationships `tree:"-"`
}

type VirtualNetworkPeeringRelationships struct {
	SourceVirtualNetwork *VirtualNetwork
	RemoteVirtualNetwork *VirtualNetwork
}
