package elasticache

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	e := &ElastiCache{
		Clusters: []Cluster{
			{
				Resource:           resource.Resource{ID: "cluster-1"},
				ParameterGroupName: value.New("pg-1", 0, "", nil),
				ReplicationGroupID: value.New("rg-1", 0, "", nil),
			},
		},
		ReplicationGroups: []ReplicationGroup{
			{
				Resource: resource.Resource{ID: "rg-1"},
				ID:       value.New("rg-1", 0, "", nil),
			},
		},
		ParameterGroups: []ParameterGroup{
			{
				Resource: resource.Resource{ID: "pg-1"},
				Name:     value.New("pg-1", 0, "", nil),
			},
		},
	}

	e.PostProcess()
	pg := e.Clusters[0].Relationships.ParameterGroup
	rg := e.Clusters[0].Relationships.ReplicationGroup

	e.PostProcess()
	assert.Equal(t, pg, e.Clusters[0].Relationships.ParameterGroup)
	assert.Equal(t, rg, e.Clusters[0].Relationships.ReplicationGroup)
	assert.NotNil(t, pg)
	assert.NotNil(t, rg)
}
