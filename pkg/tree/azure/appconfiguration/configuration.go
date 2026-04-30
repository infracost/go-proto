package appconfiguration

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Configuration struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	Replicas          value.Int    `tree:"replicas"`
}
