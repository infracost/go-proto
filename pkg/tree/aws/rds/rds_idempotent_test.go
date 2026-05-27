package rds

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	r := &RDS{
		Clusters: []Cluster{
			{
				Resource:   resource.Resource{ID: "cluster-1"},
				Identifier: value.New("cluster-1", 0, "", nil),
			},
		},
		Instances: []Instance{
			{
				Resource:                    resource.Resource{ID: "i-1"},
				Identifier:                  value.New("i-1", 0, "", nil),
				ClusterID:                   value.New("cluster-1", 0, "", nil),
			},
			{
				Resource:                    resource.Resource{ID: "i-2"},
				Identifier:                  value.New("i-2", 0, "", nil),
				ReplicateSourceDBIdentifier: value.New("i-1", 0, "", nil),
			},
		},
	}

	r.PostProcess()
	clusterInstances := append([]*Instance(nil), r.Clusters[0].Relationships.Instances...)
	i1Cluster := r.Instances[0].Relationships.Cluster
	i2Replica := r.Instances[1].Relationships.ReplicateSourceDB

	r.PostProcess()
	assert.Equal(t, clusterInstances, r.Clusters[0].Relationships.Instances)
	assert.Equal(t, i1Cluster, r.Instances[0].Relationships.Cluster)
	assert.Equal(t, i2Replica, r.Instances[1].Relationships.ReplicateSourceDB)
	assert.Len(t, clusterInstances, 1)
	assert.NotNil(t, i1Cluster)
	assert.NotNil(t, i2Replica)
}
