package route53

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ResolverFirewallRuleGroupAssociation struct {
	resource.Resource `tree:"-"`
}

type ResolverFirewallDomainList struct {
	resource.Resource `tree:"-"`
	Domains           value.Int `tree:"domains"`
}
