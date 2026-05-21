package compute

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostProcess_IGMInheritsTemplateFields(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:    resource.Resource{ID: "igm-1"},
				Name:        value.New("my-igm", 0, "", nil),
				TemplateRef: value.New("tpl-1", 0, "", nil),
				TargetSize:  value.New(int64(3), 0, "", nil),
			},
		},
		InstanceTemplates: []InstanceTemplate{
			{
				Resource:    resource.Resource{ID: "tpl-1"},
				Name:        value.New("my-template", 0, "", nil),
				MachineType: value.New("n1-standard-4", 0, "", nil),
				ScratchDisks: value.New(int64(2), 0, "", nil),
				DiskData: []DiskData{
					{StorageGB: value.New(100.0, 0, "", nil)},
				},
				GuestAccelerators: []GuestAccelerator{
					{Type: value.New("nvidia-tesla-t4", 0, "", nil), Count: value.New(int64(1), 0, "", nil)},
				},
			},
		},
	}

	c.PostProcess()

	igm := c.InstanceGroupManagers[0]
	require.NotNil(t, igm.Relationships.InstanceTemplate)
	assert.Equal(t, "n1-standard-4", igm.MachineType.Value())
	assert.Equal(t, int64(2), igm.ScratchDisks.Value())
	require.Len(t, igm.DiskData, 1)
	assert.Equal(t, 100.0, igm.DiskData[0].StorageGB.Value())
	require.Len(t, igm.GuestAccelerators, 1)
	assert.Equal(t, "nvidia-tesla-t4", igm.GuestAccelerators[0].Type.Value())
}

func TestPostProcess_IGMLinkedByName(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:    resource.Resource{ID: "igm-1"},
				TemplateRef: value.New("my-template", 0, "", nil),
			},
		},
		InstanceTemplates: []InstanceTemplate{
			{
				Resource:    resource.Resource{ID: "tpl-1"},
				Name:        value.New("my-template", 0, "", nil),
				MachineType: value.New("e2-micro", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.NotNil(t, c.InstanceGroupManagers[0].Relationships.InstanceTemplate)
	assert.Equal(t, "e2-micro", c.InstanceGroupManagers[0].MachineType.Value())
}

func TestPostProcess_PerInstanceConfigLinkedToIGM(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:   resource.Resource{ID: "igm-1"},
				Name:       value.New("my-igm", 0, "", nil),
				TargetSize: value.New(int64(3), 0, "", nil),
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

	assert.Equal(t, int64(3), c.InstanceGroupManagers[0].TargetSize.Value())
	assert.Len(t, c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs, 1)
}

func TestPostProcess_PerInstanceConfigMatchesByName(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:   resource.Resource{ID: "igm-1"},
				Name:       value.New("my-igm", 0, "", nil),
				TargetSize: value.New(int64(2), 0, "", nil),
			},
		},
		PerInstanceConfigs: []PerInstanceConfig{
			{
				Resource:                resource.Resource{ID: "pic-1"},
				InstanceGroupManagerRef: value.New("my-igm", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.Len(t, c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs, 1)
}

func TestPostProcess_MultiplePerInstanceConfigsAllLinked(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:   resource.Resource{ID: "igm-1"},
				Name:       value.New("my-igm", 0, "", nil),
				TargetSize: value.New(int64(1), 0, "", nil),
			},
		},
		PerInstanceConfigs: []PerInstanceConfig{
			{
				Resource:                resource.Resource{ID: "pic-1"},
				InstanceGroupManagerRef: value.New("igm-1", 0, "", nil),
			},
			{
				Resource:                resource.Resource{ID: "pic-2"},
				InstanceGroupManagerRef: value.New("igm-1", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.Len(t, c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs, 2)
}

// PostProcess must be idempotent — the parser runs it once on its tree, then the
// processor runs it again on a re-hydrated copy after proto serialization.
func TestPostProcess_PerInstanceConfigIdempotent(t *testing.T) {
	c := &Compute{
		InstanceGroupManagers: []InstanceGroupManager{
			{
				Resource:   resource.Resource{ID: "igm-1"},
				Name:       value.New("my-igm", 0, "", nil),
				TargetSize: value.New(int64(0), 0, "", nil),
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
	c.PostProcess()

	assert.Len(t, c.InstanceGroupManagers[0].Relationships.PerInstanceConfigs, 1)
}

func TestPostProcess_RegionDiskMarkedAsAttached(t *testing.T) {
	c := &Compute{
		Instances: []Instance{
			{
				Resource:     resource.Resource{ID: "inst-1"},
				AttachedDisk: value.New("rd-1", 0, "", nil),
			},
		},
		RegionDisks: []RegionDisk{
			{
				Resource: resource.Resource{ID: "rd-1"},
				Name:     value.New("my-region-disk", 0, "", nil),
			},
			{
				Resource: resource.Resource{ID: "rd-2"},
				Name:     value.New("other-disk", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.True(t, c.RegionDisks[0].IsAttached)
	assert.False(t, c.RegionDisks[1].IsAttached)
}

func TestPostProcess_RegionDiskMatchedByName(t *testing.T) {
	c := &Compute{
		Instances: []Instance{
			{
				Resource:     resource.Resource{ID: "inst-1"},
				AttachedDisk: value.New("my-region-disk", 0, "", nil),
			},
		},
		RegionDisks: []RegionDisk{
			{
				Resource: resource.Resource{ID: "rd-1"},
				Name:     value.New("my-region-disk", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.True(t, c.RegionDisks[0].IsAttached)
}

func TestPostProcess_RegionDiskMatchedBySelfLink(t *testing.T) {
	c := &Compute{
		Instances: []Instance{
			{
				Resource:     resource.Resource{ID: "inst-1"},
				AttachedDisk: value.New("projects/p/regions/r/disks/my-disk", 0, "", nil),
			},
		},
		RegionDisks: []RegionDisk{
			{
				Resource: resource.Resource{ID: "rd-1"},
				Name:     value.New("my-disk", 0, "", nil),
				SelfLink: value.New("projects/p/regions/r/disks/my-disk", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.True(t, c.RegionDisks[0].IsAttached)
}

func TestPostProcess_DiskStillMarkedAsAttached(t *testing.T) {
	c := &Compute{
		Instances: []Instance{
			{
				Resource:     resource.Resource{ID: "inst-1"},
				AttachedDisk: value.New("disk-1", 0, "", nil),
			},
		},
		Disks: []Disk{
			{
				Resource: resource.Resource{ID: "disk-1"},
				Name:     value.New("my-disk", 0, "", nil),
			},
		},
	}

	c.PostProcess()

	assert.True(t, c.Disks[0].IsAttached)
}
