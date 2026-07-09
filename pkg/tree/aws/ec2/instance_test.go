package ec2_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceFields(t *testing.T) {
	inst := ec2.Instance{
		Type: value.New("t3.micro", 0, "", nil),
	}

	assert.Equal(t, "t3.micro", inst.Type.Value())
}

func TestInstanceZeroValue(t *testing.T) {
	var inst ec2.Instance
	assert.Equal(t, "", string(inst.Type.Value()))
}

func TestInstanceProtoRoundTrip(t *testing.T) {
	original := ec2.Instance{
		Type: value.New("t3.micro", 0, "", nil),
		AMI:  value.New("ami-0123456789abcdef0", 0, "", nil),
	}

	obj := tree.StructToValueObject(original)

	require.Contains(t, obj.Entries, "instance_type")
	assert.Equal(t, "t3.micro", obj.Entries["instance_type"].GetStringValue())
	require.Contains(t, obj.Entries, "ami")
	assert.Equal(t, "ami-0123456789abcdef0", obj.Entries["ami"].GetStringValue())

	var restored ec2.Instance
	tree.ValueObjectToStruct(obj, &restored)

	assert.Equal(t, original.Type.Value(), restored.Type.Value())
	assert.Equal(t, original.AMI.Value(), restored.AMI.Value())
}
