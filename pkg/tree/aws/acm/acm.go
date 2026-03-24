package acm

import (
	"strings"

	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
)

type CertificateManager struct {
	Certificates []Certificate `tree:"certificate"`
}

func (cm *CertificateManager) AddCertificateAuthorities(cas *acmpca.PCACertificateAuthority) {
	// reset the certificate relationships
	for i := range cm.Certificates {
		cm.Certificates[i].Relationships = CertificateRelationships{}
	}

	// iterate over the acmpca CertificateAuthorities and add them.
	for i, cert := range cm.Certificates {
		if cert.CertificateAuthorityARN.IsEmpty() {
			continue
		}

		for _, ca := range cas.CertificateAuthorities {
			if strings.Contains(cert.CertificateAuthorityARN.Value(), ca.ID) {
				caCopy := ca
				cm.Certificates[i].Relationships.CertificateAuthority = &caCopy
			}
		}
	}
}
