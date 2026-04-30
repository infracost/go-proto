package eventhub

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Namespace struct {
	resource.Resource  `tree:"-"`
	SKU                value.Value[NamespaceSKU] `tree:"sku"`
	Capacity           value.Int    `tree:"capacity"`
	DedicatedClusterID value.String `tree:"dedicated_cluster_id"`
}
