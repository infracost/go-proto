package neptune

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	n := &Neptune{
		Clusters: []Cluster{
			{
				Resource:   resource.Resource{ID: "cluster-1"},
				Identifier: value.New("cluster-1", 0, "", nil),
			},
		},
		ClusterInstances: []ClusterInstance{
			{
				Resource:          resource.Resource{ID: "ci-1"},
				ClusterIdentifier: value.New("cluster-1", 0, "", nil),
			},
		},
		ClusterSnapshots: []ClusterSnapshot{
			{
				Resource:            resource.Resource{ID: "cs-1"},
				DBClusterIdentifier: value.New("cluster-1", 0, "", nil),
			},
		},
	}

	n.PostProcess()
	instCluster := n.ClusterInstances[0].Relationships.Cluster
	snapCluster := n.ClusterSnapshots[0].Relationships.Cluster

	n.PostProcess()
	assert.Equal(t, instCluster, n.ClusterInstances[0].Relationships.Cluster)
	assert.Equal(t, snapCluster, n.ClusterSnapshots[0].Relationships.Cluster)
	assert.NotNil(t, instCluster)
	assert.NotNil(t, snapCluster)
}
