package container

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	c := &Container{
		Clusters: []Cluster{
			{
				Resource: resource.Resource{ID: "cluster-1"},
				Name:     value.New("cluster-1", 0, "", nil),
			},
		},
		NodePools: []NodePool{
			{
				Resource:  resource.Resource{ID: "np-1", IsChild: true},
				Name:      value.New("default-pool", 0, "", nil),
				ClusterID: value.New("cluster-1", 0, "", nil),
			},
			{
				Resource:  resource.Resource{ID: "np-2"},
				Name:      value.New("my-pool", 0, "", nil),
				ClusterID: value.New("cluster-1", 0, "", nil),
			},
		},
	}

	c.PostProcess()
	nps := append([]*NodePool(nil), c.Clusters[0].Relationships.NodePools...)
	defaultNP := c.Clusters[0].Relationships.DefaultNodePool
	np1Cluster := c.NodePools[0].Relationships.Cluster
	np2Cluster := c.NodePools[1].Relationships.Cluster

	c.PostProcess()
	assert.Equal(t, nps, c.Clusters[0].Relationships.NodePools)
	assert.Equal(t, defaultNP, c.Clusters[0].Relationships.DefaultNodePool)
	assert.Equal(t, np1Cluster, c.NodePools[0].Relationships.Cluster)
	assert.Equal(t, np2Cluster, c.NodePools[1].Relationships.Cluster)

	assert.Len(t, nps, 2)
	assert.NotNil(t, defaultNP)
	assert.NotNil(t, np1Cluster)
	assert.NotNil(t, np2Cluster)
}
