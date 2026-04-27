package directconnect

type DirectConnect struct {
	Connections         []Connection         `tree:"connections"`
	GatewayAssociations []GatewayAssociation `tree:"gateway_associations"`
}
