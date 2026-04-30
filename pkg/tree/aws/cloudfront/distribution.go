package cloudfront

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Distribution struct {
	resource.Resource     `tree:"-"`
	DefaultRootObject     value.String      `tree:"default_root_object"`
	WebACLID              value.String      `tree:"web_acl_id"`
	Origins               []Origin          `tree:"origins"`
	DefaultCacheBehavior  CacheBehavior     `tree:"default_cache_behavior"`
	OrderedCacheBehaviors []CacheBehavior   `tree:"ordered_cache_behaviors"`
	LoggingConfig         LoggingConfig     `tree:"logging_config"`
	ViewerCertificate     ViewerCertificate `tree:"viewer_certificate"`
}

type ViewerProtocolPolicy uint32

const (
	ViewerProtocolPolicyUnknown ViewerProtocolPolicy = iota
	ViewerProtocolPolicyAllowAll
	ViewerProtocolPolicyHTTPSOnly
	ViewerProtocolPolicyRedirectToHTTPS
)

type SSLSupportMethod uint32

const (
	SSLSupportMethodUnknown SSLSupportMethod = iota
	SSLSupportMethodSniOnly
	SSLSupportMethodVip
	SSLSupportMethodStaticIP
)

type Origin struct {
	DomainName     value.String   `tree:"domain_name"`
	S3OriginConfig S3OriginConfig `tree:"s3_origin_config"`
	OriginShield   OriginShield   `tree:"origin_shield"`
}

type S3OriginConfig struct {
	OriginAccessIdentity value.String `tree:"origin_access_identity"`
}

type OriginShield struct {
	Enabled            value.Bool   `tree:"enabled"`
	OriginShieldRegion value.String `tree:"origin_shield_region"`
}

type CacheBehavior struct {
	ViewerProtocolPolicy   value.Value[ViewerProtocolPolicy] `tree:"viewer_protocol_policy"`
	FieldLevelEncryptionID value.String                      `tree:"field_level_encryption_id"`
}

type LoggingConfig struct {
	Bucket value.String `tree:"bucket"`
}

type ViewerCertificate struct {
	SSLSupportMethod value.Value[SSLSupportMethod] `tree:"ssl_support_method"`
}
