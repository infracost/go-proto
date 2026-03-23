package value

import (
	"testing"

	prototree "github.com/infracost/proto/gen/go/infracost/tree"
	"github.com/stretchr/testify/assert"
)

func TestStringValue(t *testing.T) {
	v := New("hello", 0, "", nil)
	assert.Equal(t, "hello", v.Value())
}

func TestBoolValue(t *testing.T) {
	v := New(true, 0, "", nil)
	assert.Equal(t, true, v.Value())
}

func TestIntValue(t *testing.T) {
	v := New(int64(42), 0, "", nil)
	assert.Equal(t, int64(42), v.Value())
}

func TestDoubleValue(t *testing.T) {
	v := New(3.14, 0, "", nil)
	assert.Equal(t, 3.14, v.Value())
}

func TestFromProtoReusesPointer(t *testing.T) {
	p := &prototree.Value{Value: &prototree.Value_StringValue{StringValue: "test"}}
	v := FromProto[string](p)
	assert.Same(t, p, v.ToProto())
}

func TestRoundTrip(t *testing.T) {
	original := New("round-trip", 0, "", nil)
	p := original.ToProto()
	restored := FromProto[string](p)
	assert.Equal(t, "round-trip", restored.Value())
	assert.Same(t, p, restored.ToProto())
}

func TestZeroValue(t *testing.T) {
	var v Value[string]
	assert.Equal(t, "", v.Value())
}
