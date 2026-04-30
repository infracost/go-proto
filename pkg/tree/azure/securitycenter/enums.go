package securitycenter

type PricingTier uint32

const (
	PricingTierUnknown PricingTier = iota
	PricingTierFree
	PricingTierStandard
)
