package appservice

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_FunctionAppTierFromServicePlan(t *testing.T) {
	cases := []struct {
		name     string
		skuName  string
		wantTier AppServiceTier
	}{
		{name: "standard", skuName: "S1", wantTier: AppServiceTierStandard},
		{name: "elastic premium", skuName: "EP1", wantTier: AppServiceTierPremium},
		{name: "elastic premium is case-insensitive", skuName: "ep2", wantTier: AppServiceTierPremium},
		{name: "flex consumption", skuName: "FC1", wantTier: AppServiceTierFlexConsumption},
		{name: "flex consumption is case-insensitive", skuName: "fc1", wantTier: AppServiceTierFlexConsumption},
		{name: "consumption", skuName: "Y1", wantTier: AppServiceTierStandard},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &AppService{
				ServicePlans: []ServicePlan{
					{
						Resource: resource.Resource{ID: "plan-1"},
						SKUName:  value.New(tc.skuName, 0, "", nil),
						OSType:   value.New("linux", 0, "", nil),
					},
				},
				FunctionApps: []FunctionApp{
					{
						Resource:         resource.Resource{ID: "func-1"},
						AppServicePlanID: value.New("plan-1", 0, "", nil),
					},
				},
			}

			s.PostProcess()

			assert.Equal(t, tc.wantTier, s.FunctionApps[0].Tier.Value())
			assert.Equal(t, tc.skuName, s.FunctionApps[0].SKU.Value())
		})
	}
}
