package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FunctionApp struct {
	resource.Resource  `tree:"-"`
	SKU                value.String `tree:"sku"`
	Tier               value.Value[AppServiceTier] `tree:"tier"`
	OSType             value.String `tree:"os_type"`
	MinTLSVersion      value.String `tree:"min_tls_version"`
	HTTPSOnly          value.Bool   `tree:"https_only"`
	AppServicePlanID   value.String `tree:"app_service_plan_id"`
	StorageAccountName value.String `tree:"storage_account_name"`
}
