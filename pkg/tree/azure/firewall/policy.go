package firewall

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Policy struct {
	resource.Resource      `tree:"-"`
	SKU                    value.Value[PolicySKU]              `tree:"sku"`
	ThreatIntelligenceMode value.Value[ThreatIntelligenceMode] `tree:"threat_intelligence_mode"`
	IntrusionDetectionMode value.Value[IntrusionDetectionMode] `tree:"intrusion_detection_mode"`
	TLSCertificateName    value.String                        `tree:"tls_certificate_name"`

	Relationships PolicyRelationships `tree:"-"`
}

type PolicyRelationships struct {
	CollectionGroups []CollectionGroup
}
