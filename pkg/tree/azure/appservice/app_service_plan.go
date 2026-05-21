package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type AppServicePlan struct {
	resource.Resource `tree:"-"`
	SKUName           value.String `tree:"sku_name"`
	SKUSize           value.String `tree:"sku_size"`
	SKUCapacity       value.Int    `tree:"sku_capacity"`
	WorkerCount       value.Int    `tree:"worker_count"`
	OSType            value.String `tree:"os_type"`
	Kind              value.Value[AppServiceKind] `tree:"kind"`
	Tier              value.Value[AppServiceTier] `tree:"tier"`
}
