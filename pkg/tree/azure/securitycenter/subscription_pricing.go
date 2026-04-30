package securitycenter

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SubscriptionPricing struct {
	resource.Resource `tree:"-"`
	Tier              value.Value[PricingTier] `tree:"tier"`
	ResourceType      value.String `tree:"resource_type"`
}
