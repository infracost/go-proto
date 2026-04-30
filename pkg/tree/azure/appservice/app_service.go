package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type App struct {
	resource.Resource `tree:"-"`
	MinTLSVersion     value.String `tree:"min_tls_version"`
	HTTPSOnly         value.Bool   `tree:"https_only"`
	AppServicePlanID  value.String `tree:"app_service_plan_id"`
}
