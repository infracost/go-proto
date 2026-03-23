package ec2

import (
	"testing"

	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostProcess_LinkInstanceStateToInstance(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{Resource: resource.Resource{ID: "i-111"}},
			{Resource: resource.Resource{ID: "i-222"}},
		},
		InstanceStates: []InstanceStateMapping{
			{InstanceID: value.New("i-222", 0, "", nil)},
		},
	}

	ec2.PostProcess()

	assert.Nil(t, ec2.Instances[0].Relationships.InstanceState)
	require.NotNil(t, ec2.Instances[1].Relationships.InstanceState)
	assert.Equal(t, &ec2.InstanceStates[0], ec2.Instances[1].Relationships.InstanceState)
}

func TestPostProcess_LinkLaunchTemplateByID(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-abc", 0, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{Resource: resource.Resource{ID: "lt-abc"}},
		},
	}

	ec2.PostProcess()

	require.NotNil(t, ec2.Instances[0].Relationships.LaunchTemplate)
	assert.Equal(t, &ec2.LaunchTemplates[0], ec2.Instances[0].Relationships.LaunchTemplate)
}

func TestPostProcess_LinkLaunchTemplateByName(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:           resource.Resource{ID: "i-111"},
				LaunchTemplateName: value.New("my-lt", 0, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{
				Resource: resource.Resource{ID: "lt-abc"},
				Name:     value.New("my-lt", 0, "", nil),
			},
		},
	}

	ec2.PostProcess()

	require.NotNil(t, ec2.Instances[0].Relationships.LaunchTemplate)
	assert.Equal(t, &ec2.LaunchTemplates[0], ec2.Instances[0].Relationships.LaunchTemplate)
}

func TestPostProcess_LaunchTemplateInheritsInstanceType(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-abc", 0, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{
				Resource:     resource.Resource{ID: "lt-abc"},
				InstanceType: value.New("t3.large", 0, "", nil),
			},
		},
	}

	ec2.PostProcess()

	assert.Equal(t, "t3.large", ec2.Instances[0].Type.Value())
}

func TestPostProcess_InstanceTypeNotOverriddenByLaunchTemplate(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-abc", 0, "", nil),
				Type:             value.New("t3.micro", 0, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{
				Resource:     resource.Resource{ID: "lt-abc"},
				InstanceType: value.New("t3.large", 0, "", nil),
			},
		},
	}

	ec2.PostProcess()

	assert.Equal(t, "t3.micro", ec2.Instances[0].Type.Value())
}

func TestPostProcess_LaunchTemplateInheritsEBSOptimized(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-abc", 0, "", nil),
				EBSOptimized:     value.New(false, flag.Defaulted, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{
				Resource:     resource.Resource{ID: "lt-abc"},
				EBSOptimized: value.New(true, 0, "", nil),
			},
		},
	}

	ec2.PostProcess()

	assert.True(t, ec2.Instances[0].EBSOptimized.IsTrue())
}

func TestPostProcess_LaunchTemplateMergesBlockDevices(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-abc", 0, "", nil),
				BlockDeviceMappings: []BlockDeviceMapping{
					{
						DeviceName: value.New("/dev/sda1", 0, "", nil),
						EBSVolume: EBSVolume{
							Size: value.New(int64(20), 0, "", nil),
						},
					},
				},
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{
				Resource: resource.Resource{ID: "lt-abc"},
				BlockDeviceMappings: []BlockDeviceMapping{
					{
						DeviceName: value.New("/dev/sda1", 0, "", nil),
						EBSVolume: EBSVolume{
							Size: value.New(int64(50), 0, "", nil),
							IOPS: value.New(int64(3000), 0, "", nil),
						},
					},
					{
						DeviceName: value.New("/dev/sdb", 0, "", nil),
						EBSVolume: EBSVolume{
							Size: value.New(int64(100), 0, "", nil),
						},
					},
				},
			},
		},
	}

	ec2.PostProcess()

	bdms := ec2.Instances[0].BlockDeviceMappings
	require.Len(t, bdms, 2)

	// find the two devices by name
	var sda1, sdb *BlockDeviceMapping
	for i := range bdms {
		switch bdms[i].DeviceName.Value() {
		case "/dev/sda1":
			sda1 = &bdms[i]
		case "/dev/sdb":
			sdb = &bdms[i]
		}
	}

	require.NotNil(t, sda1, "expected /dev/sda1")
	require.NotNil(t, sdb, "expected /dev/sdb")

	// existing device keeps its size, inherits IOPS from launch template
	assert.Equal(t, int64(20), sda1.EBSVolume.Size.Value())
	assert.Equal(t, int64(3000), sda1.EBSVolume.IOPS.Value())

	// new device added from launch template
	assert.Equal(t, int64(100), sdb.EBSVolume.Size.Value())
}

func TestPostProcess_NoLaunchTemplateMatch(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-111"},
				LaunchTemplateID: value.New("lt-xyz", 0, "", nil),
			},
		},
		LaunchTemplates: []LaunchTemplate{
			{Resource: resource.Resource{ID: "lt-abc"}},
		},
	}

	ec2.PostProcess()

	assert.Nil(t, ec2.Instances[0].Relationships.LaunchTemplate)
}

func TestPostProcess_ResetsRelationships(t *testing.T) {
	lt := LaunchTemplate{Resource: resource.Resource{ID: "lt-old"}}
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource: resource.Resource{ID: "i-111"},
				Relationships: InstanceRelationships{
					LaunchTemplate: &lt,
				},
			},
		},
	}

	ec2.PostProcess()

	assert.Nil(t, ec2.Instances[0].Relationships.LaunchTemplate)
	assert.Nil(t, ec2.Instances[0].Relationships.InstanceState)
}
