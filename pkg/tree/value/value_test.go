package value

import (
	"testing"

	"github.com/infracost/go-proto/pkg/flag"
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

func TestEnumValue(t *testing.T) {
	type MyEnum uint32
	const MyEnumVal MyEnum = 3
	v := New(MyEnumVal, 0, "", nil)
	assert.Equal(t, MyEnumVal, v.Value())
}

func TestIsDefault(t *testing.T) {
	v := New("hello", flag.Defaulted, "", nil)
	assert.True(t, v.IsDefault())

	v2 := New("hello", 0, "", nil)
	assert.False(t, v2.IsDefault())
}

func TestIsDefault_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.True(t, v.IsDefault())
}

func TestIsEmpty(t *testing.T) {
	v := New("", 0, "", nil)
	assert.True(t, v.IsEmpty())

	v2 := New("hello", 0, "", nil)
	assert.False(t, v2.IsEmpty())
}

func TestIsEmpty_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.True(t, v.IsEmpty())
}

func TestIsEmpty_Bool(t *testing.T) {
	v := New(false, 0, "", nil)
	assert.True(t, v.IsEmpty())

	v2 := New(true, 0, "", nil)
	assert.False(t, v2.IsEmpty())
}

func TestIsDefaultOrEmpty(t *testing.T) {
	// default and empty
	v := New("", flag.Defaulted, "", nil)
	assert.True(t, v.IsDefaultOrEmpty())

	// not default but empty
	v2 := New("", 0, "", nil)
	assert.True(t, v2.IsDefaultOrEmpty())

	// default but not empty
	v3 := New("hello", flag.Defaulted, "", nil)
	assert.True(t, v3.IsDefaultOrEmpty())

	// neither
	v4 := New("hello", 0, "", nil)
	assert.False(t, v4.IsDefaultOrEmpty())
}

func TestEquals(t *testing.T) {
	v := New("hello", 0, "", nil)
	assert.True(t, v.Equals("hello"))
	assert.False(t, v.Equals("world"))
}

func TestEquals_Int(t *testing.T) {
	v := New(int64(42), 0, "", nil)
	assert.True(t, v.Equals(42))
	assert.False(t, v.Equals(99))
}

func TestString(t *testing.T) {
	v := New("hello", 0, "", nil)
	assert.Equal(t, "hello", v.String())
}

func TestWithValue(t *testing.T) {
	v := New("hello", 0, "field_name", nil)
	v2 := v.WithValue("world")
	assert.Equal(t, "world", v2.Value())
	assert.Equal(t, "field_name", v2.Field())
}

func TestField(t *testing.T) {
	v := New("x", 0, "my_field", nil)
	assert.Equal(t, "my_field", v.Field())
}

func TestField_Empty(t *testing.T) {
	v := New("x", 0, "", nil)
	assert.Equal(t, "", v.Field())
}

func TestFlags(t *testing.T) {
	v := New("x", flag.Defaulted|flag.Synthetic, "", nil)
	assert.True(t, v.Flags().IsDefault())
	assert.True(t, v.Flags().IsSynthetic())
}

func TestFlags_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.Equal(t, flag.Flags(0), v.Flags())
}

func TestSource_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.Nil(t, v.Source())
}

func TestIsTrue(t *testing.T) {
	v := New(true, 0, "", nil)
	assert.True(t, v.IsTrue())
	assert.False(t, v.IsFalse())
}

func TestIsFalse(t *testing.T) {
	v := New(false, 0, "", nil)
	assert.True(t, v.IsFalse())
	assert.False(t, v.IsTrue())
}

func TestIsTrue_ZeroValue(t *testing.T) {
	var v Value[bool]
	assert.False(t, v.IsTrue())
	assert.True(t, v.IsFalse())
}

func TestIsSynthetic(t *testing.T) {
	v := New("x", flag.Synthetic, "", nil)
	assert.True(t, v.IsSynthetic())

	v2 := New("x", 0, "", nil)
	assert.False(t, v2.IsSynthetic())
}

func TestIsSynthetic_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.False(t, v.IsSynthetic())
}

func TestToProto_ZeroValue(t *testing.T) {
	var v Value[string]
	p := v.ToProto()
	assert.NotNil(t, p)
	assert.Equal(t, "", p.GetStringValue())
}

func TestSetProto(t *testing.T) {
	var v Value[string]
	p := &prototree.Value{Value: &prototree.Value_StringValue{StringValue: "injected"}}
	v.SetProto(p)
	assert.Equal(t, "injected", v.Value())
}

func TestNewList(t *testing.T) {
	items := []Value[string]{
		New("a", 0, "", nil),
		New("b", 0, "", nil),
	}
	l := NewList(items, 0, "tags", nil)
	assert.Len(t, l.Items(), 2)
	assert.Equal(t, "a", l.Items()[0].Value())
	assert.Equal(t, "b", l.Items()[1].Value())
}

func TestNewList_Empty(t *testing.T) {
	l := NewList([]Value[int64]{}, 0, "", nil)
	assert.Empty(t, l.Items())
}
