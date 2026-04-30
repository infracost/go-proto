package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Firewall struct {
	resource.Resource `tree:"-"`
	SKUTier           value.Value[FirewallSKUTier] `tree:"sku_tier"`
	VirtualHubCount   value.Int                    `tree:"virtual_hub_count"`
}
