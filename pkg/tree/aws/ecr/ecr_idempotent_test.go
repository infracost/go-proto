package ecr

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	e := &ECR{
		Repositories: []Repository{
			{
				Resource: resource.Resource{ID: "repo-1"},
				Name:     value.New("repo-1", 0, "", nil),
			},
		},
		LifecyclePolicies: []LifecyclePolicy{
			{
				Resource:       resource.Resource{ID: "lp-1"},
				RepositoryName: value.New("repo-1", 0, "", nil),
			},
		},
	}

	e.PostProcess()
	lps := append([]*LifecyclePolicy(nil), e.Repositories[0].Relationships.LifecyclePolicies...)
	assert.Len(t, lps, 1)

	e.PostProcess()
	assert.Equal(t, lps, e.Repositories[0].Relationships.LifecyclePolicies)
}
