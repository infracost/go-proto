package container

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &Container{
		KubernetesClusters: []KubernetesCluster{
			{Resource: resource.Resource{ID: "cluster-1"}},
		},
		KubernetesClusterNodePools: []KubernetesClusterNodePool{
			{
				Resource:  resource.Resource{ID: "np-1"},
				ClusterID: value.New("cluster-1", 0, "", nil),
			},
		},
		AppEnvironments: []AppEnvironment{
			{
				Resource: resource.Resource{ID: "env-1"},
				WorkloadProfiles: []WorkloadProfile{
					{Name: value.New("wp-1", 0, "", nil)},
				},
			},
		},
		Apps: []App{
			{
				Resource:                  resource.Resource{ID: "app-1"},
				ContainerAppEnvironmentID: value.New("env-1", 0, "", nil),
				WorkloadProfileName:       value.New("wp-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	npCluster := s.KubernetesClusterNodePools[0].Relationships.Cluster
	appEnv := s.Apps[0].Relationships.Environment
	appWP := s.Apps[0].Relationships.WorkloadProfile

	s.PostProcess()
	assert.Equal(t, npCluster, s.KubernetesClusterNodePools[0].Relationships.Cluster)
	assert.Equal(t, appEnv, s.Apps[0].Relationships.Environment)
	assert.Equal(t, appWP, s.Apps[0].Relationships.WorkloadProfile)
	assert.NotNil(t, npCluster)
	assert.NotNil(t, appEnv)
	assert.NotNil(t, appWP)
}
