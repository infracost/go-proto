package storage

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Storage{
		Accounts: []Account{
			{Resource: resource.Resource{ID: "acc-1"}},
		},
		ManagementPolicies: []ManagementPolicy{
			{
				Resource:         resource.Resource{ID: "mp-1"},
				StorageAccountID: value.New("acc-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	mp := s.Accounts[0].Relationships.ManagementPolicy

	s.PostProcess()
	assert.Equal(t, mp, s.Accounts[0].Relationships.ManagementPolicy)
	assert.NotNil(t, mp)
}
