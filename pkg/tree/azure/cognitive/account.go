package cognitive

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Account struct {
	resource.Resource `tree:"-"`
	Kind              value.String `tree:"kind"`
	SKU               value.String `tree:"sku"`
}
