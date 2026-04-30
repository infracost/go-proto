package globalaccelerator

type GlobalAccelerator struct {
	Accelerators   []Accelerator   `tree:"accelerators"`
	EndpointGroups []EndpointGroup `tree:"endpoint_groups"`
}
