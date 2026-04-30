package iothub

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DPS struct {
	resource.Resource `tree:"-"`
	SKU               value.Value[DPSSKU] `tree:"sku"`
}
