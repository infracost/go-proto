package route53

type Route53 struct {
	HealthChecks      []HealthCheck      `tree:"health_checks"`
	Records           []Record           `tree:"records"`
	Zones             []Zone             `tree:"zones"`
	ResolverEndpoints []ResolverEndpoint `tree:"resolver_endpoints"`

	ResolverFirewallRuleGroupAssociations []ResolverFirewallRuleGroupAssociation `tree:"resolver_firewall_rule_group_associations"`
	ResolverFirewallDomainLists           []ResolverFirewallDomainList           `tree:"resolver_firewall_domain_lists"`
	Profiles                              []Profile                              `tree:"profiles"`
	TrafficPolicyInstances                []TrafficPolicyInstance                `tree:"traffic_policy_instances"`
}

func (r *Route53) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range r.Records {
		r.Records[i].Relationships.AliasRecord = nil
	}

	// link alias records to their target records
	for i, recA := range r.Records {
		for j, recB := range r.Records {
			if i != j && !recA.AliasName.IsEmpty() && recA.AliasName.Value() == recB.Name.Value() {
				r.Records[i].Relationships.AliasRecord = &r.Records[j]
			}
		}
	}
}
