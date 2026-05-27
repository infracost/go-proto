package recoveryservices

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &RecoveryServices{
		Vaults: []Vault{
			{
				Resource: resource.Resource{ID: "vault-1"},
				Name:     value.New("vault-1", 0, "", nil),
			},
		},
		BackupProtectedVMs: []BackupProtectedVM{
			{
				Resource:          resource.Resource{ID: "bpvm-1"},
				RecoveryVaultName: value.New("vault-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	pvms := append([]*BackupProtectedVM(nil), s.Vaults[0].Relationships.ProtectedVMs...)
	assert.Len(t, pvms, 1)

	s.PostProcess()
	assert.Equal(t, pvms, s.Vaults[0].Relationships.ProtectedVMs)
}
