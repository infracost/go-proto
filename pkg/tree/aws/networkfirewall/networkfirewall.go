package networkfirewall

type NetworkFirewall struct {
	Firewalls               []Firewall               `tree:"firewalls"`
	VPCEndpointAssociations []VPCEndpointAssociation `tree:"vpc_endpoint_associations"`
}
