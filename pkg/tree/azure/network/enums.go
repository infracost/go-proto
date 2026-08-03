package network

type PublicIPAllocationMethod uint32

const (
	PublicIPAllocationMethodUnknown PublicIPAllocationMethod = iota
	PublicIPAllocationMethodStatic
	PublicIPAllocationMethodDynamic
)

type PublicIPSKU uint32

const (
	PublicIPSKUUnknown PublicIPSKU = iota
	PublicIPSKUBasic
	PublicIPSKUStandard
	PublicIPSKUStandardV2
)

type PublicIPSKUTier uint32

const (
	PublicIPSKUTierUnknown PublicIPSKUTier = iota
	PublicIPSKUTierRegional
	PublicIPSKUTierGlobal
)

type LoadBalancerSKU uint32

const (
	LoadBalancerSKUUnknown LoadBalancerSKU = iota
	LoadBalancerSKUBasic
	LoadBalancerSKUStandard
	LoadBalancerSKUGateway
)

type FirewallSKUTier uint32

const (
	FirewallSKUTierUnknown FirewallSKUTier = iota
	FirewallSKUTierStandard
	FirewallSKUTierPremium
	FirewallSKUTierBasic
)

type GatewayConnectionType uint32

const (
	GatewayConnectionTypeUnknown GatewayConnectionType = iota
	GatewayConnectionTypeIPsec
	GatewayConnectionTypeExpressRoute
	GatewayConnectionTypeVnet2Vnet
	GatewayConnectionTypeVPNClient
)

type VPNGatewayType uint32

const (
	VPNGatewayTypeUnknown VPNGatewayType = iota
	VPNGatewayTypeRouteBased
	VPNGatewayTypePolicyBased
)

type CDNOptimizationType uint32

const (
	CDNOptimizationTypeUnknown CDNOptimizationType = iota
	CDNOptimizationTypeGeneralWebDelivery
	CDNOptimizationTypeDynamicSiteAcceleration
	CDNOptimizationTypeGeneralMediaStreaming
	CDNOptimizationTypeVideoOnDemandMediaStreaming
	CDNOptimizationTypeLargeFileDownload
)
