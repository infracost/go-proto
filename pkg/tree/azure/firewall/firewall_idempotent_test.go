package firewall

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Firewall{
		Policies: []Policy{
			{Resource: resource.Resource{ID: "policy-1"}},
		},
		CollectionGroups: []CollectionGroup{
			{
				Resource:         resource.Resource{ID: "cg-1"},
				FirewallPolicyID: value.New("policy-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	cg := append([]CollectionGroup(nil), s.Policies[0].Relationships.CollectionGroups...)
	assert.Len(t, cg, 1)

	s.PostProcess()
	assert.Equal(t, cg, s.Policies[0].Relationships.CollectionGroups)
}
