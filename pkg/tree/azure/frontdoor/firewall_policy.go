package frontdoor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FirewallPolicy struct {
	resource.Resource `tree:"-"`
	CustomRules       value.Int `tree:"custom_rules"`
	ManagedRulesets   value.Int `tree:"managed_rulesets"`
}
