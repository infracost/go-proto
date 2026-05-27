package network

type Network struct {
	ApplicationGateways              []ApplicationGateway              `tree:"application_gateways"`
	BastionHosts                     []BastionHost                     `tree:"bastion_hosts"`
	ConnectionMonitors               []ConnectionMonitor               `tree:"connection_monitors"`
	DDoSProtectionPlans              []DDoSProtectionPlan              `tree:"ddos_protection_plans"`
	Firewalls                        []Firewall                        `tree:"firewalls"`
	LoadBalancers                    []LoadBalancer                    `tree:"load_balancers"`
	LoadBalancerRules                []LoadBalancerRule                `tree:"load_balancer_rules"`
	NATGateways                      []NATGateway                      `tree:"nat_gateways"`
	PrivateEndpoints                 []PrivateEndpoint                 `tree:"private_endpoints"`
	PublicIPs                        []PublicIP                        `tree:"public_ips"`
	PublicIPPrefixes                 []PublicIPPrefix                  `tree:"public_ip_prefixes"`
	TrafficManagerProfiles           []TrafficManagerProfile           `tree:"traffic_manager_profiles"`
	TrafficManagerEndpoints          []TrafficManagerEndpoint          `tree:"traffic_manager_endpoints"`
	VirtualHubs                      []VirtualHub                      `tree:"virtual_hubs"`
	VirtualNetworkGateways           []VirtualNetworkGateway           `tree:"virtual_network_gateways"`
	VirtualNetworkGatewayConnections []VirtualNetworkGatewayConnection `tree:"virtual_network_gateway_connections"`
	VirtualNetworkPeerings           []VirtualNetworkPeering           `tree:"virtual_network_peerings"`
	VirtualNetworks                  []VirtualNetwork                  `tree:"virtual_networks"`
	VPNGateways                      []VPNGateway                      `tree:"vpn_gateways"`
	VPNGatewayConnections            []VPNGatewayConnection            `tree:"vpn_gateway_connections"`
	P2PVPNGateways                   []P2PVPNGateway                   `tree:"p2p_vpn_gateways"`
	Watchers                         []Watcher                         `tree:"watchers"`
	WatcherFlowLogs                  []WatcherFlowLog                  `tree:"watcher_flow_logs"`
	FrontDoors                       []FrontDoor                       `tree:"front_doors"`
	CDNProfiles                      []CDNProfile                      `tree:"cdn_profiles"`
	CDNEndpoints                     []CDNEndpoint                     `tree:"cdn_endpoints"`
}

func (s *Network) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.LoadBalancerRules {
		s.LoadBalancerRules[i].Relationships.LoadBalancer = nil
	}
	for i := range s.VirtualNetworkGatewayConnections {
		s.VirtualNetworkGatewayConnections[i].Relationships.Gateway = nil
	}
	for i := range s.VirtualNetworkPeerings {
		s.VirtualNetworkPeerings[i].Relationships = VirtualNetworkPeeringRelationships{}
	}
	for i := range s.TrafficManagerEndpoints {
		s.TrafficManagerEndpoints[i].Relationships.Profile = nil
	}
	for i := range s.CDNEndpoints {
		s.CDNEndpoints[i].Relationships.Profile = nil
	}

	// Link LB rules to LBs
	for i, rule := range s.LoadBalancerRules {
		for j := range s.LoadBalancers {
			if rule.LoadBalancerID.Value() == s.LoadBalancers[j].ID {
				s.LoadBalancerRules[i].Relationships.LoadBalancer = &s.LoadBalancers[j]
				break
			}
		}
	}

	// Link gateway connections to gateways
	for i, connection := range s.VirtualNetworkGatewayConnections {
		for j := range s.VirtualNetworkGateways {
			if connection.VirtualNetworkGatewayID.Value() == s.VirtualNetworkGateways[j].ID {
				s.VirtualNetworkGatewayConnections[i].Relationships.Gateway = &s.VirtualNetworkGateways[j]
				break
			}
		}
	}

	// Link peerings to VNets
	for i, peering := range s.VirtualNetworkPeerings {
		for j := range s.VirtualNetworks {
			if peering.RemoteVirtualNetworkID.Value() == s.VirtualNetworks[j].ID {
				s.VirtualNetworkPeerings[i].Relationships.RemoteVirtualNetwork = &s.VirtualNetworks[j]
			}
			if peering.VirtualNetworkName.Value() == s.VirtualNetworks[j].Name.Value() {
				s.VirtualNetworkPeerings[i].Relationships.SourceVirtualNetwork = &s.VirtualNetworks[j]
			}
		}
	}

	// Link traffic manager endpoints to profiles
	for i, endpoint := range s.TrafficManagerEndpoints {
		for j := range s.TrafficManagerProfiles {
			if endpoint.ProfileID.Value() == s.TrafficManagerProfiles[j].ID {
				s.TrafficManagerEndpoints[i].Relationships.Profile = &s.TrafficManagerProfiles[j]
				break
			}
		}
	}

	// Link CDN endpoints to profiles
	for i, endpoint := range s.CDNEndpoints {
		for j := range s.CDNProfiles {
			if endpoint.ProfileName.Value() == s.CDNProfiles[j].Name.Value() {
				s.CDNEndpoints[i].Relationships.Profile = &s.CDNProfiles[j]
				break
			}
		}
	}
}
