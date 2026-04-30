package acmpca

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CertificateAuthority struct {
	resource.Resource `tree:"-"`

	UsageMode value.Value[UsageMode] `tree:"usage_mode"`
}

type UsageMode uint32

const (
	UsageModeGeneralPurpose UsageMode = iota
	UsageModeShortLivedCertificate
)
