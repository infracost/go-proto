package firewall

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CollectionGroup struct {
	resource.Resource `tree:"-"`
	FirewallPolicyID  value.String     `tree:"firewall_policy_id"`
	CollectionRules   []CollectionRule `tree:"collection_rules"`
}

type CollectionRule struct {
	TerminatedTLS   value.Bool        `tree:"terminated_tls"`
	WebCategories   value.List[string] `tree:"web_categories"`
	DestinationURLs value.List[string] `tree:"destination_urls"`
}
