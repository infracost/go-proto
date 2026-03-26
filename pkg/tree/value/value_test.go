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
	assert.True(t, v.Equal("hello"))
	assert.False(t, v.Equal("world"))
}

func TestEquals_Int(t *testing.T) {
	v := New(int64(42), 0, "", nil)
	assert.True(t, v.Equal(42))
	assert.False(t, v.Equal(99))
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

func TestContains(t *testing.T) {
	v := New("hello world", 0, "", nil)
	assert.True(t, v.Contains("world"))
	assert.True(t, v.Contains("hello"))
	assert.False(t, v.Contains("missing"))
}

func TestContains_NonString(t *testing.T) {
	v := New(int64(42), 0, "", nil)
	assert.False(t, v.Contains("42"))
}

func TestHasPrefix(t *testing.T) {
	v := New("hello world", 0, "", nil)
	assert.True(t, v.HasPrefix("hello"))
	assert.False(t, v.HasPrefix("world"))
}

func TestHasSuffix(t *testing.T) {
	v := New("hello world", 0, "", nil)
	assert.True(t, v.HasSuffix("world"))
	assert.False(t, v.HasSuffix("hello"))
}

func TestIsEqualFold(t *testing.T) {
	v := New("Hello", 0, "", nil)
	assert.True(t, v.EqualFold("hello"))
	assert.True(t, v.EqualFold("HELLO"))
	assert.False(t, v.EqualFold("world"))
}

func TestIsOneOf(t *testing.T) {
	v := New("b", 0, "", nil)
	assert.True(t, v.IsOneOf("a", "b", "c"))
	assert.False(t, v.IsOneOf("x", "y", "z"))
}

func TestIsOneOf_Int(t *testing.T) {
	v := New(int64(2), 0, "", nil)
	assert.True(t, v.IsOneOf(1, 2, 3))
	assert.False(t, v.IsOneOf(4, 5, 6))
}

func TestValueOr(t *testing.T) {
	v := New("hello", 0, "", nil)
	assert.Equal(t, "hello", v.ValueOr("fallback"))

	empty := New("", 0, "", nil)
	assert.Equal(t, "fallback", empty.ValueOr("fallback"))
}

func TestValueOr_ZeroValue(t *testing.T) {
	var v Value[string]
	assert.Equal(t, "default", v.ValueOr("default"))
}

func TestValueOr_Int(t *testing.T) {
	v := New(int64(0), 0, "", nil)
	assert.Equal(t, int64(99), v.ValueOr(99))

	v2 := New(int64(42), 0, "", nil)
	assert.Equal(t, int64(42), v2.ValueOr(99))
}

func TestIsGreaterThan_Int(t *testing.T) {
	v := New(int64(10), 0, "", nil)
	assert.True(t, v.IsGreaterThan(5))
	assert.False(t, v.IsGreaterThan(10))
	assert.False(t, v.IsGreaterThan(15))
}

func TestIsLessThan_Int(t *testing.T) {
	v := New(int64(10), 0, "", nil)
	assert.True(t, v.IsLessThan(15))
	assert.False(t, v.IsLessThan(10))
	assert.False(t, v.IsLessThan(5))
}

func TestIsGreaterThanOrEqual_Int(t *testing.T) {
	v := New(int64(10), 0, "", nil)
	assert.True(t, v.IsGreaterThanOrEqual(5))
	assert.True(t, v.IsGreaterThanOrEqual(10))
	assert.False(t, v.IsGreaterThanOrEqual(15))
}

func TestIsLessThanOrEqual_Int(t *testing.T) {
	v := New(int64(10), 0, "", nil)
	assert.True(t, v.IsLessThanOrEqual(15))
	assert.True(t, v.IsLessThanOrEqual(10))
	assert.False(t, v.IsLessThanOrEqual(5))
}

func TestIsGreaterThan_Float(t *testing.T) {
	v := New(3.14, 0, "", nil)
	assert.True(t, v.IsGreaterThan(2.0))
	assert.False(t, v.IsGreaterThan(3.14))
	assert.False(t, v.IsGreaterThan(4.0))
}

func TestIsGreaterThan_String(t *testing.T) {
	v := New("banana", 0, "", nil)
	assert.True(t, v.IsGreaterThan("apple"))
	assert.False(t, v.IsGreaterThan("banana"))
	assert.False(t, v.IsGreaterThan("cherry"))
}

func TestIsLessThan_String(t *testing.T) {
	v := New("banana", 0, "", nil)
	assert.True(t, v.IsLessThan("cherry"))
	assert.False(t, v.IsLessThan("banana"))
	assert.False(t, v.IsLessThan("apple"))
}

func TestComparison_Bool(t *testing.T) {
	// bool comparisons always return 0 (equal)
	v := New(true, 0, "", nil)
	assert.False(t, v.IsGreaterThan(false))
	assert.False(t, v.IsLessThan(false))
}

func TestPointer(t *testing.T) {
	v := New("hello", 0, "", nil)
	p := v.Pointer()
	assert.Equal(t, "hello", *p)
}

func TestPointer_Int(t *testing.T) {
	v := New(int64(42), 0, "", nil)
	p := v.Pointer()
	assert.Equal(t, int64(42), *p)
}
