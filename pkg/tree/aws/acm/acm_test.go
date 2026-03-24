package acm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestCertificateManager_AddCertificateAuthorities(t *testing.T) {
	tests := []struct {
		name         string
		certificates []Certificate
		cas          []acmpca.CertificateAuthority
		wantLinkedCA []string // expected linked CA ID per cert, empty string means no link
	}{
		{
			name:         "no certificates or CAs",
			certificates: nil,
			cas:          nil,
			wantLinkedCA: nil,
		},
		{
			name: "certificate with empty ARN is not linked",
			certificates: []Certificate{
				{CertificateAuthorityARN: value.EmptyString},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []string{""},
		},
		{
			name: "certificate with matching ARN is linked",
			certificates: []Certificate{
				{CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123", 0, "", nil)},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []string{"arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"},
		},
		{
			name: "certificate with non-matching ARN is not linked",
			certificates: []Certificate{
				{CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/xyz-999", 0, "", nil)},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []string{""},
		},
		{
			name: "multiple certificates with mixed matching",
			certificates: []Certificate{
				{CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123", 0, "", nil)},
				{CertificateAuthorityARN: value.EmptyString},
				{CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/def-456", 0, "", nil)},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/def-456"}},
			},
			wantLinkedCA: []string{
				"arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123",
				"",
				"arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/def-456",
			},
		},
		{
			name: "existing relationships are reset",
			certificates: []Certificate{
				{
					CertificateAuthorityARN: value.EmptyString,
					Relationships: CertificateRelationships{
						CertificateAuthority: &acmpca.CertificateAuthority{
							Resource: resource.Resource{ID: "stale-ca"},
						},
					},
				},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []string{""},
		},
		{
			name: "linked CA is a copy not a reference to the slice element",
			certificates: []Certificate{
				{CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123", 0, "", nil)},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []string{"arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CertificateManager{Certificates: tt.certificates}
			pca := &acmpca.PCACertificateAuthority{CertificateAuthorities: tt.cas}

			cm.AddCertificateAuthorities(pca)

			require.Len(t, cm.Certificates, len(tt.wantLinkedCA))
			for i, wantID := range tt.wantLinkedCA {
				cert := cm.Certificates[i]
				if wantID == "" {
					assert.Nil(t, cert.Relationships.CertificateAuthority, "cert[%d] should not have a linked CA", i)
				} else {
					require.NotNil(t, cert.Relationships.CertificateAuthority, "cert[%d] should have a linked CA", i)
					assert.Equal(t, wantID, cert.Relationships.CertificateAuthority.ID, "cert[%d] linked to wrong CA", i)
				}
			}
		})
	}
}
