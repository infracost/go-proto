package fabric

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Capacity struct {
	resource.Resource `tree:"-"`
	SKUName           value.String `tree:"sku_name"`
}
