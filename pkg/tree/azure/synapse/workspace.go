package synapse

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workspace struct {
	resource.Resource                `tree:"-"`
	ManagedVirtualNetworkEnabled     value.Bool `tree:"managed_virtual_network_enabled"`
}
