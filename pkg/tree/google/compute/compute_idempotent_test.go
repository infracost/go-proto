package compute

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:    resource.Resource{ID: "igm-1"},
				Name:        value.New("igm-1", 0, "", nil),
				TemplateRef: value.New("tpl-1", 0, "", nil),
			},
		},
		InstanceTemplates: []InstanceTemplate{
			{
				Resource: resource.Resource{ID: "tpl-1"},
				Name:     value.New("tpl-1", 0, "", nil),
			},
		},
		PerInstanceConfigs: []PerInstanceConfig{
			{
				Resource:                resource.Resource{ID: "pic-1"},
				InstanceGroupManagerRef: value.New("igm-1", 0, "", nil),
			},
		},
	}

	c.PostProcess()
	pics := append([]*PerInstanceConfig(nil), c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs...)
	tpl := c.InstanceGroupManagers[0].Relationships.InstanceTemplate

	c.PostProcess()
	assert.Equal(t, pics, c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs)
	assert.Equal(t, tpl, c.InstanceGroupManagers[0].Relationships.InstanceTemplate)

	assert.Len(t, pics, 1)
	assert.NotNil(t, tpl)
}
