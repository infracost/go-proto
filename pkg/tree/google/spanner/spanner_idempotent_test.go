package spanner

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Spanner{
		Instances: []Instance{
			{
				Resource: resource.Resource{ID: "inst-1"},
				Name:     value.New("inst-1", 0, "", nil),
			},
		},
		Databases: []Database{
			{
				Resource:     resource.Resource{ID: "db-1"},
				InstanceName: value.New("inst-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	inst := s.Databases[0].Relationships.Instance

	s.PostProcess()
	assert.Equal(t, inst, s.Databases[0].Relationships.Instance)
	assert.NotNil(t, inst)
}
