package vpclattice

type VPCLattice struct {
	Services                           []Service                           `tree:"services"`
	ServiceNetworkResourceAssociations []ServiceNetworkResourceAssociation `tree:"service_network_resource_associations"`
}
