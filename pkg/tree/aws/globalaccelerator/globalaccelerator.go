package globalaccelerator

type GlobalAccelerator struct {
	Accelerators              []Accelerator   `tree:"accelerators"`
	CustomRoutingAccelerators []Accelerator   `tree:"custom_routing_accelerators"`
	EndpointGroups            []EndpointGroup `tree:"endpoint_groups"`
}
