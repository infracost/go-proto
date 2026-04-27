package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Connection struct {
	resource.Resource    `tree:"-"`
	Bandwidth            value.String                `tree:"bandwidth"`
	Location             value.String                `tree:"location"`
	ConnectionType       value.Value[ConnectionType] `tree:"connection_type"`
	VirtualInterfaceType value.String                `tree:"virtual_interface_type"`
}

type ConnectionType uint32

const (
	ConnectionTypeUnknown   ConnectionType = 0
	ConnectionTypeDedicated ConnectionType = 1
	ConnectionTypeHosted    ConnectionType = 2
)
