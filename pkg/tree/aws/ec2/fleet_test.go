package ec2

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFleet_LaunchTemplateLink(t *testing.T) {
	e := &EC2{
		LaunchTemplates: []LaunchTemplate{
			{Resource: resource.Resource{ID: "lt-1"}, Name: value.New("byname", 0, "", nil), InstanceType: value.New("m5.large", 0, "", nil)},
		},
		Fleets: []Fleet{
			{LaunchTemplateID: value.New("lt-1", 0, "", nil)},
			{LaunchTemplateName: value.New("byname", 0, "", nil)},
			{LaunchTemplateID: value.New("lt-1", 0, "", nil), InstanceType: value.New("c5.large", 0, "", nil)},
			{},
		},
	}

	e.PostProcess()
	e.PostProcess()

	require.NotNil(t, e.Fleets[0].Relationships.LaunchTemplate)
	assert.Equal(t, "m5.large", e.Fleets[0].EffectiveInstanceType())
	assert.Equal(t, "m5.large", e.Fleets[1].EffectiveInstanceType())
	assert.Equal(t, "c5.large", e.Fleets[2].EffectiveInstanceType())
	assert.Equal(t, "", e.Fleets[3].EffectiveInstanceType())
}
