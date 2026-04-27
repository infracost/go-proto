package cloudfront

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Distribution struct {
	resource.Resource      `tree:"-"`
	DefaultRootObject      value.String      `tree:"default_root_object"`
	WebACLID               value.String      `tree:"web_acl_id"`
	Origins                []Origin          `tree:"origins"`
	DefaultCacheBehaviour  CacheBehaviour    `tree:"default_cache_behaviour"`
	OrderedCacheBehaviours []CacheBehaviour  `tree:"ordered_cache_behaviours"`
	LoggingConfig          *LoggingConfig    `tree:"logging_config"`
	ViwerCertificate       ViewerCertificate `tree:"viewer_certificate"`
}

type ViewerProtocolPolicy uint32

const (
	ViewerProtocolPolicyAllowAll ViewerProtocolPolicy = iota
	ViewerProtocolPolicyHTTPSOnly
	ViewerProtocolPolicyRedirectToHTTPS
)

type SSLSupportMethod uint32

const (
	SSLSupportMethodSniOnly SSLSupportMethod = iota
	SSLSupportMethodVip
	SSLSupportMethodStaticIP
)

type Origin struct {
	DomainName     value.String    `tree:"domain_name"`
	S3OriginConfig *S3OriginConfig `tree:"s3_origin_config"`
	OriginShield   *OriginShield   `tree:"origin_shield"`
}

type S3OriginConfig struct {
	OriginAccessIdentity value.String `tree:"origin_access_identity"`
}

type OriginShield struct {
	Enabled            value.Bool   `tree:"enabled"`
	OriginShieldRegion value.String `tree:"origin_shield_region"`
}

type CacheBehaviour struct {
	ViewerProtocolPolicy   value.Value[ViewerProtocolPolicy] `tree:"viewer_protocol_policy"`
	FieldLevelEncryptionID value.String                      `tree:"field_level_encryption_id"`
}

type LoggingConfig struct {
	Bucket value.String `tree:"bucket"`
}

type ViewerCertificate struct {
	SSLSupportMethod value.Value[SSLSupportMethod] `tree:"ssl_support_method"`
}
