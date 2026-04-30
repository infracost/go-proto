package appservice

type AppService struct {
	Apps                   []App                   `tree:"app_services"`
	CertificateBindings    []CertificateBinding    `tree:"certificate_bindings"`
	Certificates           []Certificate           `tree:"certificates"`
	CertificateOrders      []CertificateOrder      `tree:"certificate_orders"`
	CustomHostnameBindings []CustomHostnameBinding  `tree:"custom_hostname_bindings"`
	Environments           []Environment            `tree:"environments"`
	AppServicePlans        []AppServicePlan         `tree:"app_service_plans"`
	ServicePlans           []ServicePlan            `tree:"service_plans"`
	FunctionApps           []FunctionApp            `tree:"function_apps"`
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
				s.FunctionApps[i].OSType = plan.Kind
				if !plan.Tier.IsEmpty() {
					s.FunctionApps[i].Tier = plan.Tier
				}
				break
			}
		}

		for _, plan := range s.ServicePlans {
			if plan.ID == fa.AppServicePlanID.Value() {
				s.FunctionApps[i].SKU = plan.SKUName
				s.FunctionApps[i].OSType = plan.OSType
				break
			}
		}
	}
}
