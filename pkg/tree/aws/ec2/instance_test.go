package ec2

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/convert"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceFields(t *testing.T) {
	inst := Instance{
		Type: value.New("t3.micro", 0, "", nil),
	}

	assert.Equal(t, "t3.micro", inst.Type.Value())
}

func TestInstanceZeroValue(t *testing.T) {
	var inst Instance
	assert.Equal(t, "", string(inst.Type.Value()))
}

func TestInstanceProtoRoundTrip(t *testing.T) {
	original := Instance{
		Type: value.New("t3.micro", 0, "", nil),
	}

	obj := convert.StructToValueObject(original)

	require.Contains(t, obj.Entries, "instance_type")
	assert.Equal(t, "t3.micro", obj.Entries["instance_type"].GetStringValue())

	var restored Instance
	convert.ValueObjectToStruct(obj, &restored)

	assert.Equal(t, original.Type.Value(), restored.Type.Value())
}
