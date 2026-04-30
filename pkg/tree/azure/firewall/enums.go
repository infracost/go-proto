package firewall

type PolicySKU uint32

const (
	PolicySKUUnknown PolicySKU = iota
	PolicySKUBasic
	PolicySKUStandard
	PolicySKUPremium
)

type ThreatIntelligenceMode uint32

const (
	ThreatIntelligenceModeUnknown ThreatIntelligenceMode = iota
	ThreatIntelligenceModeOff
	ThreatIntelligenceModeAlert
	ThreatIntelligenceModeDeny
)

type IntrusionDetectionMode uint32

const (
	IntrusionDetectionModeUnknown IntrusionDetectionMode = iota
	IntrusionDetectionModeOff
	IntrusionDetectionModeAlert
	IntrusionDetectionModeDeny
)
