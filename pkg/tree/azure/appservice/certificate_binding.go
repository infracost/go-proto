package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CertificateBinding struct {
	resource.Resource `tree:"-"`
	SSLState          value.Value[SSLState] `tree:"ssl_state"`
	CertificateID     value.String `tree:"certificate_id"`

	Relationships CertificateBindingRelationships `tree:"-"`
}

type CertificateBindingRelationships struct {
	Certificate *Certificate
}
