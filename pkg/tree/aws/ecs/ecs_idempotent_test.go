package ecs

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	e := &ECS{
		Clusters: []Cluster{
			{
				Resource: resource.Resource{ID: "cluster-1"},
				Name:     value.New("cluster-1", 0, "", nil),
			},
		},
		ClusterCapacityProviders: []ClusterCapacityProviders{
			{
				Resource:    resource.Resource{ID: "ccp-1"},
				ClusterName: value.New("cluster-1", 0, "", nil),
			},
		},
		CapacityProviders: []CapacityProvider{
			{
				Resource:             resource.Resource{ID: "cp-1"},
				Name:                 value.New("FARGATE", 0, "", nil),
				CapacityProviderType: value.New(CapacityProviderTypeFargate, 0, "", nil),
			},
		},
		TaskDefinitions: []TaskDefinition{
			{Resource: resource.Resource{ID: "td-1"}},
		},
		TaskSets: []TaskSet{
			{
				Resource:                resource.Resource{ID: "ts-1"},
				TaskDefinitionReference: value.New("td-1", 0, "", nil),
				ServiceReference:        value.New("svc-1", 0, "", nil),
			},
		},
		Services: []Service{
			{
				Resource:                resource.Resource{ID: "svc-1"},
				ClusterReference:        value.New("cluster-1", 0, "", nil),
				TaskDefinitionReference: value.New("td-1", 0, "", nil),
			},
		},
	}

	e.PostProcess()
	cps := append([]*ClusterCapacityProviders(nil), e.Clusters[0].Relationships.CapacityProviders...)
	taskSetTD := e.TaskSets[0].Relationships.TaskDefinition
	svcCluster := e.Services[0].Relationships.Cluster
	svcTD := e.Services[0].Relationships.TaskDefinition

	e.PostProcess()
	assert.Equal(t, cps, e.Clusters[0].Relationships.CapacityProviders)
	assert.Equal(t, taskSetTD, e.TaskSets[0].Relationships.TaskDefinition)
	assert.Equal(t, svcCluster, e.Services[0].Relationships.Cluster)
	assert.Equal(t, svcTD, e.Services[0].Relationships.TaskDefinition)

	assert.Len(t, cps, 1)
	assert.NotNil(t, taskSetTD)
	assert.NotNil(t, svcCluster)
	assert.NotNil(t, svcTD)
}
