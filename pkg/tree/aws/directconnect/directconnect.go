package directconnect

type DirectConnect struct {
	Connections         []Connection         `tree:"connections"`
	GatewayAssociations []GatewayAssociation `tree:"gateway_associations"`
	HostedConnections   []HostedConnection   `tree:"hosted_connections"`
	VirtualInterfaces   []VirtualInterface   `tree:"virtual_interfaces"`
}
