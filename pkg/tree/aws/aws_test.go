package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestAWS_PostProcess_ACM(t *testing.T) {
	tests := []struct {
		name         string
		certificates []acm.Certificate
		cas          []acmpca.CertificateAuthority
		wantLinkedCA []bool // whether each cert should have a linked CA
	}{
		{
			name:         "no certificates or CAs",
			certificates: nil,
			cas:          nil,
			wantLinkedCA: nil,
		},
		{
			name: "certificate with empty ARN is not linked",
			certificates: []acm.Certificate{
				{CertificateAuthorityARN: value.EmptyString},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []bool{false},
		},
		{
			name: "certificate with matching ARN is linked",
			certificates: []acm.Certificate{
				{
					CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123", 0, "", nil),
				},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []bool{true},
		},
		{
			name: "certificate with non-matching ARN is not linked",
			certificates: []acm.Certificate{
				{
					CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/xyz-999", 0, "", nil),
				},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []bool{false},
		},
		{
			name: "multiple certificates with mixed matching",
			certificates: []acm.Certificate{
				{
					CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123", 0, "", nil),
				},
				{
					CertificateAuthorityARN: value.EmptyString,
				},
				{
					CertificateAuthorityARN: value.New("arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/def-456", 0, "", nil),
				},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/def-456"}},
			},
			wantLinkedCA: []bool{true, false, true},
		},
		{
			name: "existing relationships are reset before linking",
			certificates: []acm.Certificate{
				{
					CertificateAuthorityARN: value.EmptyString,
					Relationships: acm.CertificateRelationships{
						CertificateAuthority: &acmpca.CertificateAuthority{
							Resource: resource.Resource{ID: "stale-ca"},
						},
					},
				},
			},
			cas: []acmpca.CertificateAuthority{
				{Resource: resource.Resource{ID: "arn:aws:acm-pca:us-east-1:123456789012:certificate-authority/abc-123"}},
			},
			wantLinkedCA: []bool{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				CertificateManager: acm.CertificateManager{
					Certificates: tt.certificates,
				},
				PCACertificateAuthority: acmpca.PCACertificateAuthority{
					CertificateAuthorities: tt.cas,
				},
			}

			aws.PostProcess()

			require.Len(t, aws.CertificateManager.Certificates, len(tt.wantLinkedCA))
			for i, wantLinked := range tt.wantLinkedCA {
				cert := aws.CertificateManager.Certificates[i]
				if wantLinked {
					assert.NotNil(t, cert.Relationships.CertificateAuthority, "cert[%d] should have a linked CA", i)
				} else {
					assert.Nil(t, cert.Relationships.CertificateAuthority, "cert[%d] should not have a linked CA", i)
				}
			}
		})
	}
}
