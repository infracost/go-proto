package route53

type Route53 struct {
	HealthChecks      []HealthCheck      `tree:"health_checks"`
	Records           []Record           `tree:"records"`
	Zones             []Zone             `tree:"zones"`
	ResolverEndpoints []ResolverEndpoint `tree:"resolver_endpoints"`
}

func (r *Route53) PostProcess() {
	// link alias records to their target records
	for i, recA := range r.Records {
		for j, recB := range r.Records {
			if i != j && !recA.AliasName.IsEmpty() && recA.AliasName.Value() == recB.Name.Value() {
				r.Records[i].Relationships.AliasRecord = &r.Records[j]
			}
		}
	}
}
