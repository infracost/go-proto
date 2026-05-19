package appservice

import (
	"strings"

	"github.com/infracost/go-proto/pkg/tree/value"
)

type AppService struct {
	Apps                   []App                   `tree:"app_services"`
	CertificateBindings    []CertificateBinding    `tree:"certificate_bindings"`
	Certificates           []Certificate           `tree:"certificates"`
	CertificateOrders      []CertificateOrder      `tree:"certificate_orders"`
	CustomHostnameBindings []CustomHostnameBinding `tree:"custom_hostname_bindings"`
	Environments           []Environment           `tree:"environments"`
	AppServicePlans        []AppServicePlan        `tree:"app_service_plans"`
	ServicePlans           []ServicePlan           `tree:"service_plans"`
	FunctionApps           []FunctionApp           `tree:"function_apps"`
}

func (s *AppService) PostProcess() {
	// link certificate bindings to certificates
	for i, binding := range s.CertificateBindings {
		for j := range s.Certificates {
			if binding.CertificateID.Value() == s.Certificates[j].ID {
				s.CertificateBindings[i].Relationships.Certificate = &s.Certificates[j]
				break
			}
		}
	}

	// link function apps to app service plans and service plans
	for i, fa := range s.FunctionApps {
		if fa.AppServicePlanID.IsEmpty() {
			continue
		}

		for _, plan := range s.AppServicePlans {
			if plan.ID == fa.AppServicePlanID.Value() {
				s.FunctionApps[i].SKU = plan.SKUSize
				// Legacy behaviour: AppServicePlan.Kind drives FunctionApp.OSType.
				s.FunctionApps[i].OSType = appServiceKindToString(plan.Kind.Value())
				tier := plan.Tier.Value()
				if tier == AppServiceTierUnknown {
					tier = AppServiceTierStandard
				}
				// Function apps on elastic premium plans (kind=elastic, tier=ElasticPremium)
				// bill compute via the plan itself, so the function app is free. The
				// Y1 SKU (Consumption) is always billed per-execution regardless of
				// other fields.
				if !strings.EqualFold(plan.SKUSize.Value(), "y1") &&
					(plan.Kind.Value() == AppServiceKindElastic ||
						plan.Tier.Value() == AppServiceTierElasticPremium) {
					tier = AppServiceTierPremium
				}
				s.FunctionApps[i].Tier = value.New(tier, 0, "", nil)
				break
			}
		}

		for _, plan := range s.ServicePlans {
			if plan.ID == fa.AppServicePlanID.Value() {
				s.FunctionApps[i].SKU = plan.SKUName
				s.FunctionApps[i].OSType = plan.OSType
				tier := AppServiceTierStandard
				if strings.HasPrefix(strings.ToLower(plan.SKUName.Value()), "ep") {
					tier = AppServiceTierPremium
				}
				s.FunctionApps[i].Tier = value.New(tier, 0, "", nil)
				break
			}
		}
	}
}

func appServiceKindToString(k AppServiceKind) value.String {
	var s string
	switch k {
	case AppServiceKindApp, AppServiceKindWindows:
		s = "windows"
	case AppServiceKindLinux:
		s = "linux"
	case AppServiceKindElastic:
		s = "elastic"
	case AppServiceKindFunctionApp:
		s = "functionapp"
	}
	return value.New(s, 0, "", nil)
}
