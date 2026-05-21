package waf

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type WebACLv2 struct {
	resource.Resource `tree:"-"`
	Rules             value.Int `tree:"rules"`
	RuleGroups        value.Int `tree:"rule_groups"`
	ManagedRuleGroups value.Int `tree:"managed_rule_groups"`
}
