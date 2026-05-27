package synapse

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Synapse{
		Workspaces: []Workspace{
			{Resource: resource.Resource{ID: "ws-1"}},
		},
		SparkPools: []SparkPool{
			{
				Resource:    resource.Resource{ID: "sp-1"},
				WorkspaceID: value.New("ws-1", 0, "", nil),
			},
		},
		SQLPools: []SQLPool{
			{
				Resource:    resource.Resource{ID: "sql-1"},
				WorkspaceID: value.New("ws-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	sparkWS := s.SparkPools[0].Relationships.Workspace
	sqlWS := s.SQLPools[0].Relationships.Workspace

	s.PostProcess()
	assert.Equal(t, sparkWS, s.SparkPools[0].Relationships.Workspace)
	assert.Equal(t, sqlWS, s.SQLPools[0].Relationships.Workspace)
	assert.NotNil(t, sparkWS)
	assert.NotNil(t, sqlWS)
}
