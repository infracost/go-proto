package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type PublicIP struct {
	resource.Resource `tree:"-"`
	SKU               value.Value[PublicIPSKU]              `tree:"sku"`
	SKUTier           value.Value[PublicIPSKUTier]          `tree:"sku_tier"`
	AllocationMethod  value.Value[PublicIPAllocationMethod]  `tree:"allocation_method"`
}
