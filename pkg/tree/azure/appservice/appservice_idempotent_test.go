package appservice

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &AppService{
		Certificates: []Certificate{
			{Resource: resource.Resource{ID: "cert-1"}},
		},
		CertificateBindings: []CertificateBinding{
			{
				Resource:      resource.Resource{ID: "cb-1"},
				CertificateID: value.New("cert-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	cert := s.CertificateBindings[0].Relationships.Certificate

	s.PostProcess()
	assert.Equal(t, cert, s.CertificateBindings[0].Relationships.Certificate)
	assert.NotNil(t, cert)
}
