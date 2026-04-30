package signalr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ServiceResource struct {
	resource.Resource `tree:"-"`
	SKUName           value.String `tree:"sku_name"`
	SKUCapacity       value.Int    `tree:"sku_capacity"`
}
