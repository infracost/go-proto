package sagemaker

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostProcess_LinksEndpointConfig(t *testing.T) {
	s := &SageMaker{
		EndpointConfigurations: []EndpointConfiguration{
			{Resource: resource.Resource{ID: "cfg-1"}},
			{Resource: resource.Resource{ID: "cfg-2"}, Name: value.New("by-name", 0, "", nil)},
		},
		Endpoints: []Endpoint{
			{EndpointConfigName: value.New("cfg-1", 0, "", nil)},
			{EndpointConfigName: value.New("by-name", 0, "", nil)},
			{EndpointConfigName: value.New("missing", 0, "", nil)},
		},
	}

	s.PostProcess()
	s.PostProcess()
	require.NotNil(t, s.Endpoints[0].Relationships.EndpointConfiguration)
	assert.Equal(t, "cfg-1", s.Endpoints[0].Relationships.EndpointConfiguration.ID)
	require.NotNil(t, s.Endpoints[1].Relationships.EndpointConfiguration)
	assert.Equal(t, "cfg-2", s.Endpoints[1].Relationships.EndpointConfiguration.ID)
	assert.Nil(t, s.Endpoints[2].Relationships.EndpointConfiguration)
}
