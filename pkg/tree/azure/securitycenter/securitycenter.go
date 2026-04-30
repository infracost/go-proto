package securitycenter

type SecurityCenter struct {
	SubscriptionPricings []SubscriptionPricing `tree:"subscription_pricings"`
}
