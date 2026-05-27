package cognitive

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Cognitive{
		Accounts: []Account{
			{Resource: resource.Resource{ID: "acc-1"}},
		},
		Deployments: []Deployment{
			{
				Resource:           resource.Resource{ID: "dep-1"},
				CognitiveAccountID: value.New("acc-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	acc := s.Deployments[0].Relationships.Account

	s.PostProcess()
	assert.Equal(t, acc, s.Deployments[0].Relationships.Account)
	assert.NotNil(t, acc)
}
