package event

import (
	"testing"

	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/provider"
	"github.com/stretchr/testify/assert"
)

func makeResource(supported bool, free bool, components ...*provider.CostComponent) *provider.Resource {
	return &provider.Resource{
		IsSupported: supported,
		IsFree:      free,
		Costs: &provider.ResourceCosts{
			Components: components,
		},
	}
}

func makeMonthlyCostComponent(price, qty float64) *provider.CostComponent {
	return &provider.CostComponent{
		PeriodPrice: &provider.PeriodPrice{
			Price:  rat.New(price).Proto(),
			Period: provider.Period_MONTH,
		},
		Quantity: rat.New(qty).Proto(),
	}
}

func makeHourlyCostComponent(price, qty float64) *provider.CostComponent {
	return &provider.CostComponent{
		PeriodPrice: &provider.PeriodPrice{
			Price:  rat.New(price).Proto(),
			Period: provider.Period_HOUR,
		},
		Quantity: rat.New(qty).Proto(),
	}
}

func TestResourceMonthlyCost_SupportedResource(t *testing.T) {
	r := makeResource(true, false,
		makeMonthlyCostComponent(10, 2),  // $20
		makeMonthlyCostComponent(5, 3),   // $15
	)

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.Equals(rat.New(35)))
}

func TestResourceMonthlyCost_FreeResource(t *testing.T) {
	r := makeResource(true, true,
		makeMonthlyCostComponent(10, 2),
	)

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.IsZero())
}

func TestResourceMonthlyCost_UnsupportedResource(t *testing.T) {
	r := makeResource(false, false,
		makeMonthlyCostComponent(10, 2),
	)

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.IsZero())
}

func TestResourceMonthlyCost_HourlyPricing(t *testing.T) {
	r := makeResource(true, false,
		makeHourlyCostComponent(2, 1), // $2/hr * 730 hrs = $1460
	)

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.Equals(rat.New(1460)))
}

func TestResourceMonthlyCost_WithChildResources(t *testing.T) {
	r := &provider.Resource{
		IsSupported: true,
		IsFree:      false,
		Costs: &provider.ResourceCosts{
			Components: []*provider.CostComponent{
				makeMonthlyCostComponent(10, 1), // $10
			},
		},
		ChildResources: []*provider.Resource{
			makeResource(true, false,
				makeMonthlyCostComponent(5, 1), // $5
			),
		},
	}

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.Equals(rat.New(15)))
}

func TestResourceMonthlyCost_WithDiscount(t *testing.T) {
	discountRate, _ := rat.NewFromString("1/5") // 20% discount
	r := makeResource(true, false)
	r.Costs.Components = []*provider.CostComponent{
		{
			PeriodPrice: &provider.PeriodPrice{
				Price:  rat.New(100).Proto(),
				Period: provider.Period_MONTH,
			},
			Quantity:     rat.New(1).Proto(),
			DiscountRate: discountRate.Proto(),
		},
	}

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.Equals(rat.New(80)))
}

func TestTotalMonthlyCost(t *testing.T) {
	resources := []*provider.Resource{
		makeResource(true, false, makeMonthlyCostComponent(10, 1)),  // $10
		makeResource(true, false, makeMonthlyCostComponent(20, 1)),  // $20
		makeResource(false, false, makeMonthlyCostComponent(30, 1)), // unsupported, $0
	}

	total := TotalMonthlyCost(resources)
	assert.True(t, total.Equals(rat.New(30)))
}

func TestProjectCostInfoFromResources(t *testing.T) {
	current := []*provider.Resource{
		makeResource(true, false, makeMonthlyCostComponent(100, 1)),
	}
	past := []*provider.Resource{
		makeResource(true, false, makeMonthlyCostComponent(80, 1)),
	}

	info := ProjectCostInfoFromResources("my-project", current, past)
	assert.Equal(t, "my-project", info.ProjectName)
	assert.True(t, info.TotalMonthlyCost.Equals(rat.New(100)))
	assert.True(t, info.PastTotalMonthlyCost.Equals(rat.New(80)))
}

func TestResourceMonthlyCost_NilResource(t *testing.T) {
	cost := ResourceMonthlyCost(nil)
	assert.True(t, cost.IsZero())
}

func TestResourceMonthlyCost_NoCosts(t *testing.T) {
	r := &provider.Resource{
		IsSupported: true,
		IsFree:      false,
	}

	cost := ResourceMonthlyCost(r)
	assert.True(t, cost.IsZero())
}