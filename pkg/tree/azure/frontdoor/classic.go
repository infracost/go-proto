package frontdoor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Classic represents the deprecated azurerm_frontdoor resource (Front Door
// classic). The modern azurerm_cdn_frontdoor_* family uses different SKU
// pricing and is modeled separately.
type Classic struct {
	resource.Resource `tree:"-"`
	FrontendHosts     value.Int `tree:"frontend_hosts"`
	RoutingRules      value.Int `tree:"routing_rules"`
}
