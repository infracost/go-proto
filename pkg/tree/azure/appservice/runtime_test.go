package appservice_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/azure/appservice"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionAppRuntimeProtoRoundTrip(t *testing.T) {
	original := appservice.FunctionApp{
		RuntimeName:    value.New("python", 0, "", nil),
		RuntimeVersion: value.New("3.11", 0, "", nil),
	}

	obj := tree.StructToValueObject(original)

	require.Contains(t, obj.Entries, "runtime_name")
	require.Contains(t, obj.Entries, "runtime_version")
	assert.Equal(t, "python", obj.Entries["runtime_name"].GetStringValue())
	assert.Equal(t, "3.11", obj.Entries["runtime_version"].GetStringValue())

	var restored appservice.FunctionApp
	tree.ValueObjectToStruct(obj, &restored)

	assert.Equal(t, original.RuntimeName.Value(), restored.RuntimeName.Value())
	assert.Equal(t, original.RuntimeVersion.Value(), restored.RuntimeVersion.Value())
}
