package ecs_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceProtoRoundTrip(t *testing.T) {
	original := ecs.Service{
		Name:            value.New("my-svc", 0, "", nil),
		PlatformVersion: value.New("1.4.0", 0, "", nil),
	}

	obj := tree.StructToValueObject(original)

	require.Contains(t, obj.Entries, "platform_version")
	assert.Equal(t, "1.4.0", obj.Entries["platform_version"].GetStringValue())

	var restored ecs.Service
	tree.ValueObjectToStruct(obj, &restored)

	assert.Equal(t, original.PlatformVersion.Value(), restored.PlatformVersion.Value())
}
