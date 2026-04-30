package logicapps

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Standard struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	AppServicePlanID  value.String `tree:"app_service_plan_id"`
}
