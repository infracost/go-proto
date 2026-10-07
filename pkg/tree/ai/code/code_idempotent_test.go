package code

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	c := &Code{
		Clients: []Client{
			{Resource: resource.Resource{ID: "client-1"}},
		},
		ModelConfigs: []ModelConfig{
			{Resource: resource.Resource{ID: "config-1"}},
		},
		ModelCalls: []ModelCall{
			{
				Resource:      resource.Resource{ID: "call-1"},
				ClientID:      value.New("client-1", 0, "", nil),
				ModelConfigID: value.New("config-1", 0, "", nil),
			},
			{
				// Neither reference found: both relationships stay nil.
				Resource: resource.Resource{ID: "call-2"},
			},
		},
	}

	c.PostProcess()
	assert.Same(t, &c.Clients[0], c.ModelCalls[0].Relationships.Client)
	assert.Same(t, &c.ModelConfigs[0], c.ModelCalls[0].Relationships.ModelConfig)
	assert.Nil(t, c.ModelCalls[1].Relationships.Client)
	assert.Nil(t, c.ModelCalls[1].Relationships.ModelConfig)

	c.PostProcess()
	assert.Same(t, &c.Clients[0], c.ModelCalls[0].Relationships.Client)
	assert.Same(t, &c.ModelConfigs[0], c.ModelCalls[0].Relationships.ModelConfig)
	assert.Nil(t, c.ModelCalls[1].Relationships.Client)
	assert.Nil(t, c.ModelCalls[1].Relationships.ModelConfig)
}
