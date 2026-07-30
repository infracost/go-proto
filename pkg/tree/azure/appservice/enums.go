package appservice

// AppServiceTier represents the SKU tier of an App Service Plan, and indirectly
// the billing model for any Function App or App Service hosted on it.
//
// The set of tiers is stable and curated by Azure; new tiers ship infrequently
// and require enum updates so that downstream pricing code can react.
type AppServiceTier uint32

const (
	AppServiceTierUnknown AppServiceTier = iota
	AppServiceTierFree
	AppServiceTierShared
	AppServiceTierBasic
	AppServiceTierStandard
	AppServiceTierPremium
	AppServiceTierPremiumV2
	AppServiceTierPremiumV3
	AppServiceTierElasticPremium
	AppServiceTierIsolated
	AppServiceTierIsolatedV2
	// Dynamic is the Consumption plan tier (used by the Y1 SKU).
	AppServiceTierDynamic
	// WorkflowStandard is used by Logic Apps Standard plans (WS1/WS2/WS3).
	AppServiceTierWorkflowStandard
	// FlexConsumption is the Flex Consumption plan tier (used by the FC1 SKU).
	AppServiceTierFlexConsumption
)

// AppServiceKind reflects the App Service Plan "kind" attribute. It is a
// closed set per the AzureRM provider: app/Windows/Linux/elastic/FunctionApp.
// The legacy code paths often used Kind interchangeably with the underlying
// OS, so callers should be explicit about which they mean.
type AppServiceKind uint32

const (
	AppServiceKindUnknown AppServiceKind = iota
	AppServiceKindApp
	AppServiceKindWindows
	AppServiceKindLinux
	AppServiceKindElastic
	AppServiceKindFunctionApp
)

type SSLState uint32

const (
	SSLStateUnknown SSLState = iota
	SSLStateIPBasedEnabled
	SSLStateSniEnabled
)

type CertificateProductType uint32

const (
	CertificateProductTypeUnknown CertificateProductType = iota
	CertificateProductTypeStandard
	CertificateProductTypeWildCard
)
