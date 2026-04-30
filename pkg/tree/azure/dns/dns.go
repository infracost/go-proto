package dns

type DNS struct {
	PrivateResolverForwardingRulesets []PrivateResolverForwardingRuleset `tree:"private_resolver_forwarding_rulesets"`
	PrivateResolverInboundEndpoints  []PrivateResolverInboundEndpoint   `tree:"private_resolver_inbound_endpoints"`
	PrivateResolverOutboundEndpoints []PrivateResolverOutboundEndpoint  `tree:"private_resolver_outbound_endpoints"`
	Records                          []Record                           `tree:"records"`
	Zones                            []Zone                             `tree:"zones"`
}
