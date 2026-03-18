package event

import (
	"github.com/infracost/go-proto/pkg/rat"
	"github.com/infracost/proto/gen/go/infracost/provider"
)

var hoursInMonth = rat.New(730)

// ResourceMonthlyCost computes the total monthly cost for a single proto
// resource by summing its cost components and child resources. Only supported,
// non-free resources contribute to the total.
//
// Ported from runner: pkg/infracost/scan.go convertProtoResourceToSchema()
func ResourceMonthlyCost(resource *provider.Resource) *rat.Rat {
	if resource == nil || !resource.IsSupported || resource.IsFree {
		return rat.Zero
	}

	total := rat.Zero

	if resource.Costs != nil {
		for _, cc := range resource.Costs.Components {
			total = total.Add(costComponentMonthlyCost(cc))
		}
	}

	for _, child := range resource.ChildResources {
		total = total.Add(ResourceMonthlyCost(child))
	}

	return total
}

// TotalMonthlyCost computes the total monthly cost across a set of provider
// resources.
//
// Ported from runner: pkg/infracost/output.go calculateTotalCosts()
func TotalMonthlyCost(resources []*provider.Resource) *rat.Rat {
	total := rat.Zero
	for _, r := range resources {
		total = total.Add(ResourceMonthlyCost(r))
	}
	return total
}

// ProjectCostInfoFromResources builds a ProjectCostInfo from the current and
// past provider resource lists for a given project.
//
// The dashboard API stores this data per ProjectResult row; see
// dashboard: api/src/services/guardrails.ts lookupProjectInfos()
func ProjectCostInfoFromResources(projectName string, resources, pastResources []*provider.Resource) ProjectCostInfo {
	return ProjectCostInfo{
		ProjectName:          projectName,
		TotalMonthlyCost:     TotalMonthlyCost(resources),
		PastTotalMonthlyCost: TotalMonthlyCost(pastResources),
	}
}

// costComponentMonthlyCost computes the monthly cost for a single cost component.
//
// Ported from runner: pkg/infracost/scan.go convertProtoCostComponentToSchema()
func costComponentMonthlyCost(cc *provider.CostComponent) *rat.Rat {
	if cc.PeriodPrice == nil || cc.Quantity == nil || cc.PeriodPrice.Price == nil {
		return rat.Zero
	}

	price := applyDiscount(rat.FromProto(cc.PeriodPrice.Price), rat.FromProto(cc.DiscountRate))
	_, monthlyQty := convertQuantityByPeriod(rat.FromProto(cc.Quantity), cc.PeriodPrice.Period)

	return monthlyQty.Mul(price)
}

// convertQuantityByPeriod converts a quantity to both hourly and monthly values.
//
// Ported from runner: pkg/infracost/scan.go convertQuantityByPeriod()
func convertQuantityByPeriod(qty *rat.Rat, period provider.Period) (hourly, monthly *rat.Rat) {
	switch period {
	case provider.Period_MONTH:
		return qty.Div(hoursInMonth), qty
	case provider.Period_HOUR:
		return qty, qty.Mul(hoursInMonth)
	default:
		return rat.Zero, rat.Zero
	}
}

// applyDiscount applies a discount rate to a price.
//
// Ported from runner: pkg/infracost/scan.go applyDiscount()
func applyDiscount(price *rat.Rat, discountRate *rat.Rat) *rat.Rat {
	if discountRate != nil && discountRate.GreaterThanZero() {
		return price.Mul(rat.New(1).Sub(discountRate))
	}
	return price
}