package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ServicePlan struct {
	resource.Resource `tree:"-"`
	SKUName           value.String `tree:"sku_name"`
	WorkerCount       value.Int    `tree:"worker_count"`
	OSType            value.String `tree:"os_type"`
}
