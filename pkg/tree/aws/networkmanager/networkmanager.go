package networkmanager

type NetworkManager struct {
	CoreNetworks           []CoreNetwork           `tree:"core_networks"`
	Attachments            []Attachment            `tree:"attachments"`
	TransitGatewayPeerings []TransitGatewayPeering `tree:"transit_gateway_peerings"`
}
