package acmpca

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CertificateAuthority struct {
	resource.Resource `tree:"-"`

	UsageMode       value.String
	MonthlyRequests value.Int
}
