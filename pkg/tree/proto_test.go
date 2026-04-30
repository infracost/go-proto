package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToProtoAttributes(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Type: value.New("t3.micro", 0, "", nil)},
					{Type: value.New("m5.large", 0, "", nil)},
				},
			},
		},
	}

	proto, err := tree.ToProto()
	require.NoError(t, err)

	resources := proto.Providers["aws"].Services["ec2"].Resources
	assert.Len(t, resources, 2)
	assert.Equal(t, string("t3.micro"), resources[0].Attributes.Entries["instance_type"].GetStringValue())
	assert.Equal(t, string("m5.large"), resources[1].Attributes.Entries["instance_type"].GetStringValue())
}

func TestFromProtoAttributes(t *testing.T) {
	proto := &prototree.Tree{
		Providers: map[string]*prototree.Provider{
			"aws": {
				Services: map[string]*prototree.Service{
					"ec2": {
						Resources: []*prototree.Resource{
							{
								Type: "instances",
								Attributes: &prototree.ValueObject{
									Entries: map[string]*prototree.Value{
										"instance_type": {Value: &prototree.Value_StringValue{StringValue: string("t3.micro")}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tree, err := FromProto(proto)
	require.NoError(t, err)

	require.Len(t, tree.AWS.EC2.Instances, 1)
	assert.Equal(t, "t3.micro", tree.AWS.EC2.Instances[0].Type.Value())
}

func TestRoundTrip(t *testing.T) {
	original := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Type: value.New("t3.micro", 0, "", nil)},
					{Type: value.New("m5.large", 0, "", nil)},
					{Type: value.New("c5.xlarge", 0, "", nil)},
				},
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	result, err := FromProto(proto)
	require.NoError(t, err)

	require.Len(t, result.AWS.EC2.Instances, 3)
	assert.Equal(t, "t3.micro", result.AWS.EC2.Instances[0].Type.Value())
	assert.Equal(t, "m5.large", result.AWS.EC2.Instances[1].Type.Value())
	assert.Equal(t, "c5.xlarge", result.AWS.EC2.Instances[2].Type.Value())
}

func TestFromProtoUnmappedProvider(t *testing.T) {
	proto := &prototree.Tree{
		Providers: map[string]*prototree.Provider{
			"gcp": {
				Services: map[string]*prototree.Service{},
			},
		},
	}

	_, err := FromProto(proto)
	assert.ErrorContains(t, err, "unmapped provider: gcp")
}

func TestBaseResourceRoundTrip(t *testing.T) {
	original := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{
						Resource: resource.Resource{
							ID:     "i-1234",
							Region: "us-east-1",
							Flags:  42,
						},
						Type: value.New("t3.micro", 0, "", nil),
					},
				},
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	res := proto.Providers["aws"].Services["ec2"].Resources[0]
	assert.Equal(t, "i-1234", res.Id)
	assert.Equal(t, "us-east-1", res.Region)
	assert.Equal(t, uint64(42), res.Flags)

	result, err := FromProto(proto)
	require.NoError(t, err)

	inst := result.AWS.EC2.Instances[0]
	assert.Equal(t, "i-1234", inst.ID)
	assert.Equal(t, "us-east-1", inst.Region)
	assert.Equal(t, uint64(42), inst.Flags)
	assert.Equal(t, "t3.micro", inst.Type.Value())
}

func TestUnsupportedResourcesRoundTrip(t *testing.T) {
	original := &Tree{
		UnsupportedResources: []*resource.Resource{
			{
				ID:     "unsupported-1",
				Region: "us-west-2",
				Flags:  99,
			},
			{
				ID:     "unsupported-2",
				Region: "eu-west-1",
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)
	require.Len(t, proto.UnsupportedResources, 2)
	assert.Equal(t, "unsupported-1", proto.UnsupportedResources[0].Id)
	assert.Equal(t, "us-west-2", proto.UnsupportedResources[0].Region)
	assert.Equal(t, uint64(99), proto.UnsupportedResources[0].Flags)
	assert.Equal(t, "unsupported-2", proto.UnsupportedResources[1].Id)

	result, err := FromProto(proto)
	require.NoError(t, err)
	require.Len(t, result.UnsupportedResources, 2)
	assert.Equal(t, "unsupported-1", result.UnsupportedResources[0].ID)
	assert.Equal(t, "us-west-2", result.UnsupportedResources[0].Region)
	assert.Equal(t, uint64(99), result.UnsupportedResources[0].Flags)
	assert.Equal(t, "unsupported-2", result.UnsupportedResources[1].ID)
	assert.Equal(t, "eu-west-1", result.UnsupportedResources[1].Region)
}

func TestMixedResourcesRoundTrip(t *testing.T) {
	original := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{
						Resource: resource.Resource{ID: "i-1234", Region: "us-east-1"},
						Type:     value.New("t3.micro", 0, "", nil),
					},
				},
			},
		},
		UnsupportedResources: []*resource.Resource{
			{ID: "unsupported-1", Region: "us-west-2"},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	result, err := FromProto(proto)
	require.NoError(t, err)

	require.Len(t, result.AWS.EC2.Instances, 1)
	assert.Equal(t, "i-1234", result.AWS.EC2.Instances[0].ID)
	assert.Equal(t, "t3.micro", result.AWS.EC2.Instances[0].Type.Value())

	require.Len(t, result.UnsupportedResources, 1)
	assert.Equal(t, "unsupported-1", result.UnsupportedResources[0].ID)
	assert.Equal(t, "us-west-2", result.UnsupportedResources[0].Region)
}

func TestNestedSubResourceRoundTrip(t *testing.T) {
	original := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{
						Resource: resource.Resource{
							ID:     "i-1234",
							Region: "us-east-1",
						},
						Type: value.New("t3.micro", 0, "", nil),
						RootBlockDevice: ec2.BlockDeviceMapping{
							DeviceName: value.New("/dev/sda1", 0, "", nil),
							EBSVolume: ec2.EBSVolume{
								Resource: resource.Resource{
									ID:     "vol-root",
									Region: "us-east-1",
								},
								Type:   value.New(ec2.EBSVolumeTypeGP3, 0, "", nil),
								SizeGB: value.New[int64](50, 0, "", nil),
							},
						},
						BlockDeviceMappings: []ec2.BlockDeviceMapping{
							{
								DeviceName: value.New("/dev/sdb", 0, "", nil),
								EBSVolume: ec2.EBSVolume{
									Resource: resource.Resource{
										ID:     "vol-extra",
										Region: "eu-west-1",
									},
									Type:   value.New(ec2.EBSVolumeTypeIO1, 0, "", nil),
									SizeGB: value.New[int64](100, 0, "", nil),
									IOPS:   value.New[int64](3000, 0, "", nil),
								},
							},
						},
					},
				},
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	result, err := FromProto(proto)
	require.NoError(t, err)

	require.Len(t, result.AWS.EC2.Instances, 1)
	inst := result.AWS.EC2.Instances[0]

	// top-level resource
	assert.Equal(t, "i-1234", inst.ID)
	assert.Equal(t, "us-east-1", inst.Region)

	// root block device sub-resource keeps its region
	assert.Equal(t, "vol-root", inst.RootBlockDevice.EBSVolume.ID)
	assert.Equal(t, "us-east-1", inst.RootBlockDevice.EBSVolume.Region)
	assert.Equal(t, "/dev/sda1", inst.RootBlockDevice.DeviceName.Value())
	assert.Equal(t, ec2.EBSVolumeTypeGP3, inst.RootBlockDevice.EBSVolume.Type.Value())
	assert.Equal(t, int64(50), inst.RootBlockDevice.EBSVolume.SizeGB.Value())

	// additional block device sub-resource keeps its own region
	require.Len(t, inst.BlockDeviceMappings, 1)
	extra := inst.BlockDeviceMappings[0]
	assert.Equal(t, "vol-extra", extra.EBSVolume.ID)
	assert.Equal(t, "eu-west-1", extra.EBSVolume.Region)
	assert.Equal(t, "/dev/sdb", extra.DeviceName.Value())
	assert.Equal(t, ec2.EBSVolumeTypeIO1, extra.EBSVolume.Type.Value())
	assert.Equal(t, int64(100), extra.EBSVolume.SizeGB.Value())
	assert.Equal(t, int64(3000), extra.EBSVolume.IOPS.Value())
}

func TestEmptyTree(t *testing.T) {
	tree := &Tree{}

	proto, err := tree.ToProto()
	require.NoError(t, err)
	assert.Empty(t, proto.Providers["aws"].Services["ec2"].Resources)

	result, err := FromProto(proto)
	require.NoError(t, err)
	assert.Empty(t, result.AWS.EC2.Instances)
}
