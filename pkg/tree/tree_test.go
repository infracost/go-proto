package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToResources(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
					{Resource: resource.Resource{ID: "i-222"}},
				},
			},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	inst0, ok := resources[0].(*ec2.Instance)
	require.True(t, ok)
	assert.Equal(t, "i-111", inst0.ID)

	inst1, ok := resources[1].(*ec2.Instance)
	require.True(t, ok)
	assert.Equal(t, "i-222", inst1.ID)
}

func TestToResourcesEmpty(t *testing.T) {
	tree := &Tree{}
	resources := tree.ToResources()
	assert.Empty(t, resources)
}

func TestToResourcesUnsupported(t *testing.T) {
	tree := &Tree{
		UnsupportedResources: []*resource.Resource{
			{ID: "unsupported-1", Region: "us-west-2"},
			{ID: "unsupported-2", Region: "eu-west-1"},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	assert.Equal(t, "unsupported-1", resources[0].GetBase().ID)
	assert.Equal(t, "us-west-2", resources[0].GetBase().Region)
	assert.Equal(t, "unsupported-2", resources[1].GetBase().ID)
	assert.Equal(t, "eu-west-1", resources[1].GetBase().Region)
}

func TestPostProcess_CallsChildPostProcess(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
				},
				InstanceStates: []ec2.InstanceStateMapping{
					{InstanceID: value.New("i-111", 0, "", nil)},
				},
			},
		},
	}

	tree.PostProcess()

	// EC2.PostProcess should have linked the instance state
	require.NotNil(t, tree.AWS.EC2.Instances[0].Relationships.InstanceState)
	assert.Equal(t, &tree.AWS.EC2.InstanceStates[0], tree.AWS.EC2.Instances[0].Relationships.InstanceState)
}

func TestPostProcess_EmptyTree(t *testing.T) {
	tree := &Tree{}
	// should not panic
	tree.PostProcess()
}

func TestToResourcesMixed(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
				},
			},
		},
		UnsupportedResources: []*resource.Resource{
			{ID: "unsupported-1"},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	assert.Equal(t, "i-111", resources[0].GetBase().ID)
	assert.Equal(t, "unsupported-1", resources[1].GetBase().ID)
}
