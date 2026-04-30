package databricks

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workspace struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
}
