package iothub

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Hub struct {
	resource.Resource `tree:"-"`
	SKU               value.Value[HubSKU] `tree:"sku"`
	Capacity          value.Int    `tree:"capacity"`
}
