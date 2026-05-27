package secretmanager

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &SecretManager{
		Secrets: []Secret{
			{
				Resource: resource.Resource{ID: "secret-1"},
				Name:     value.New("secret-1", 0, "", nil),
			},
		},
		SecretVersions: []SecretVersion{
			{
				Resource:  resource.Resource{ID: "sv-1"},
				SecretRef: value.New("secret-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	secret := s.SecretVersions[0].Relationships.Secret

	s.PostProcess()
	assert.Equal(t, secret, s.SecretVersions[0].Relationships.Secret)
	assert.NotNil(t, secret)
}
