package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CDNProfile struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	SKU               value.String `tree:"sku"`
}
