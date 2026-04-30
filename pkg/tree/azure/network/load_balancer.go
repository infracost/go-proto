package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LoadBalancer struct {
	resource.Resource `tree:"-"`
	SKU               value.Value[LoadBalancerSKU] `tree:"sku"`
}
