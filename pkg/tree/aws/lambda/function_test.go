package lambda_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/aws/lambda"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionProtoRoundTrip(t *testing.T) {
	original := lambda.Function{
		Name:    value.New("my-fn", 0, "", nil),
		Runtime: value.New("python3.12", 0, "", nil),
	}

	obj := tree.StructToValueObject(original)

	require.Contains(t, obj.Entries, "runtime")
	assert.Equal(t, "python3.12", obj.Entries["runtime"].GetStringValue())

	var restored lambda.Function
	tree.ValueObjectToStruct(obj, &restored)

	assert.Equal(t, original.Runtime.Value(), restored.Runtime.Value())
}
