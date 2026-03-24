package acm

import (
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Certificate struct {
	resource.Resource       `tree:"-"`
	CertificateAuthorityARN value.String `tree:"certificate_authority_arn"`

	Relationships CertificateRelationships `tree:"-"`
}

type CertificateRelationships struct {
	CertificateAuthority *acmpca.CertificateAuthority
}
