package address

import (
	"testing"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
)

func Test_Address_Parse(t *testing.T) {

	tests := []struct {
		expected *Address
		raw      string
	}{
		{
			raw: "launch_template.0.id",
			expected: &Address{
				proto: &parserpb.Address{
					Segments: []*parserpb.Segment{
						{Value: "launch_template"},
						{IndexInt: int64Ptr(0)},
						{Value: "id"},
					},
				},
			},
		},
		{
			raw: "module.x[0].resource.y",
			expected: &Address{
				proto: &parserpb.Address{
					Segments: []*parserpb.Segment{
						{Value: "module"},
						{Value: "x", IndexInt: int64Ptr(0)},
						{Value: "resource"},
						{Value: "y"},
					},
				},
			},
		},
		{
			raw: `aws_s3_bucket.this["z"]`,
			expected: &Address{
				proto: &parserpb.Address{
					Segments: []*parserpb.Segment{
						{Value: "aws_s3_bucket"},
						{Value: "this", IndexString: stringPtr("z")},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			actual, err := Parse(test.raw)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}
			if !actual.Equal(test.expected) {
				t.Errorf("\n\nexpected:\n%#v\n\ngot:\n%#v", test.expected, actual)
				t.Errorf("\n\nexpected:\n%s\n\ngot:\n%s", test.expected.String(), actual.String())
			}
		})
	}

}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func Test_Relative(t *testing.T) {

	tests := []struct {
		base   *Address
		addr   *Address
		result *Address
	}{
		{
			base:   New("a"),
			addr:   New("a", "b"),
			result: New("b"),
		},
		{
			base:   New("a"),
			addr:   New("x"),
			result: New("x"),
		},
		{
			base:   New(),
			addr:   New("a"),
			result: New("a"),
		},
		{
			base:   New("a", "b", "c"),
			addr:   New("a", "b", "c").CreateStringIndexedChild("d"),
			result: New().CreateStringIndexedChild("d"),
		},
	}

	for _, test := range tests {
		t.Run(test.addr.String(), func(t *testing.T) {
			assert.Equal(t, test.result.String(), test.addr.Relative(test.base).String())
		})
	}

}

func TestAddress_ToProto(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Nil(t, a.ToProto())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("module", "vpc")
		proto := a.ToProto()
		assert.NotNil(t, proto)
		assert.Len(t, proto.Segments, 2)
	})
}

func TestFromProto(t *testing.T) {
	t.Run("nil proto", func(t *testing.T) {
		a := FromProto(nil)
		assert.Equal(t, Empty, a)
	})

	t.Run("valid proto", func(t *testing.T) {
		proto := &parserpb.Address{
			Segments: []*parserpb.Segment{
				{Value: "test"},
			},
		}
		a := FromProto(proto)
		assert.Equal(t, "test", a.String())
	})
}

func TestAddress_MarshalJSON(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		data, err := a.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("test")
		data, err := a.MarshalJSON()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestAddress_UnmarshalJSON(t *testing.T) {
	a := New("original")
	data, _ := a.MarshalJSON()

	var b Address
	err := b.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, a.String(), b.String())
}

func TestAddress_Module(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.Module())
	})

	t.Run("empty address", func(t *testing.T) {
		a := New()
		assert.Equal(t, Empty, a.Module())
	})

	t.Run("with module", func(t *testing.T) {
		a := New("module", "vpc", "aws_subnet", "public")
		mod := a.Module()
		assert.Equal(t, "module.vpc", mod.String())
	})

	t.Run("nested modules", func(t *testing.T) {
		a := New("module", "a", "module", "b", "resource", "c")
		mod := a.Module()
		assert.Equal(t, "module.a.module.b", mod.String())
	})

	t.Run("no module", func(t *testing.T) {
		a := New("aws_instance", "test")
		mod := a.Module()
		assert.Equal(t, "", mod.String())
	})
}

func TestAddress_Local(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.Local())
	})

	t.Run("with module", func(t *testing.T) {
		a := New("module", "vpc", "aws_subnet", "public")
		local := a.Local()
		assert.Equal(t, "aws_subnet.public", local.String())
	})
}

func TestAddress_Segments(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Nil(t, a.Segments())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("a", "b")
		segs := a.Segments()
		assert.Len(t, segs, 2)
	})
}

func TestFromSegments(t *testing.T) {
	segs := []*parserpb.Segment{
		{Value: "a"},
		{Value: "b"},
	}
	a := FromSegments(segs)
	assert.Equal(t, "a.b", a.String())
}

func TestAddress_Len(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, 0, a.Len())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("a", "b", "c")
		assert.Equal(t, 3, a.Len())
	})
}

func TestAddress_IsEmpty(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.True(t, a.IsEmpty())
	})

	t.Run("empty address", func(t *testing.T) {
		a := New()
		assert.True(t, a.IsEmpty())
	})

	t.Run("non-empty address", func(t *testing.T) {
		a := New("test")
		assert.False(t, a.IsEmpty())
	})
}

func TestAddress_At(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, "", a.At(0))
	})

	t.Run("out of bounds", func(t *testing.T) {
		a := New("a")
		assert.Equal(t, "", a.At(5))
	})

	t.Run("valid index", func(t *testing.T) {
		a := New("a", "b", "c")
		assert.Equal(t, "b", a.At(1))
	})
}

func TestAddress_From(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.From(0))
	})

	t.Run("out of bounds", func(t *testing.T) {
		a := New("a")
		result := a.From(5)
		assert.True(t, result.IsEmpty())
	})

	t.Run("valid index", func(t *testing.T) {
		a := New("a", "b", "c")
		result := a.From(1)
		assert.Equal(t, "b.c", result.String())
	})
}

func TestAddress_Relative_NilCases(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.Relative(New("base")))
	})

	t.Run("nil base", func(t *testing.T) {
		a := New("a", "b")
		assert.Equal(t, a, a.Relative(nil))
	})
}

func TestAddress_Parse_Errors(t *testing.T) {
	t.Run("unclosed brackets", func(t *testing.T) {
		_, err := Parse("test[0")
		assert.Error(t, err)
	})

	t.Run("unexpected close bracket", func(t *testing.T) {
		_, err := Parse("test]")
		assert.Error(t, err)
	})

	t.Run("nested brackets", func(t *testing.T) {
		_, err := Parse("test[[0]]")
		assert.Error(t, err)
	})

	t.Run("unclosed quotes", func(t *testing.T) {
		_, err := Parse(`test["unclosed`)
		assert.Error(t, err)
	})

	t.Run("invalid index", func(t *testing.T) {
		_, err := Parse("test[abc]")
		assert.Error(t, err)
	})
}

func TestAddress_Append(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		result := a.Append("test")
		assert.Equal(t, "test", result.String())
	})

	t.Run("append string", func(t *testing.T) {
		a := New("a")
		result := a.Append("b")
		assert.Equal(t, "a.b", result.String())
	})

	t.Run("append int", func(t *testing.T) {
		a := New("a")
		result := a.Append(0)
		assert.Equal(t, "a[0]", result.String())
	})

	t.Run("append int64", func(t *testing.T) {
		a := New("a")
		result := a.Append(int64(1))
		assert.Equal(t, "a[1]", result.String())
	})

	t.Run("append float64", func(t *testing.T) {
		a := New("a")
		result := a.Append(float64(2))
		assert.Equal(t, "a[2]", result.String())
	})

	t.Run("append string slice", func(t *testing.T) {
		a := New("a")
		result := a.Append([]string{"b", "c"})
		assert.Equal(t, "a.b.c", result.String())
	})

	t.Run("append address", func(t *testing.T) {
		a := New("a")
		b := New("b", "c")
		result := a.Append(b)
		assert.Equal(t, "a.b.c", result.String())
	})

	t.Run("append proto address", func(t *testing.T) {
		a := New("a")
		proto := &parserpb.Address{Segments: []*parserpb.Segment{{Value: "b"}}}
		result := a.Append(proto)
		assert.Equal(t, "a.b", result.String())
	})

	t.Run("append nil proto", func(t *testing.T) {
		a := New("a")
		var proto *parserpb.Address
		result := a.Append(proto)
		assert.Equal(t, "a", result.String())
	})

	t.Run("append nil address", func(t *testing.T) {
		a := New("a")
		var b *Address
		result := a.Append(b)
		assert.Equal(t, "a", result.String())
	})

	t.Run("append unknown type", func(t *testing.T) {
		a := New("a")
		result := a.Append(struct{ Name string }{Name: "test"})
		assert.Contains(t, result.String(), "test")
	})
}

func TestAddress_Equal(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		var a, b *Address
		assert.True(t, a.Equal(b))
	})

	t.Run("both empty", func(t *testing.T) {
		a := New()
		b := New()
		assert.True(t, a.Equal(b))
	})

	t.Run("one nil", func(t *testing.T) {
		a := New("test")
		var b *Address
		assert.False(t, a.Equal(b))
	})

	t.Run("different lengths", func(t *testing.T) {
		a := New("a", "b")
		b := New("a")
		assert.False(t, a.Equal(b))
	})

	t.Run("same addresses", func(t *testing.T) {
		a := New("a", "b")
		b := New("a", "b")
		assert.True(t, a.Equal(b))
	})

	t.Run("different values", func(t *testing.T) {
		a := New("a", "b")
		b := New("a", "c")
		assert.False(t, a.Equal(b))
	})
}

func TestAddress_Clone(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.Clone())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("a", "b")
		clone := a.Clone()
		assert.True(t, a.Equal(clone))
		// Ensure it's a deep copy
		assert.NotSame(t, a.proto, clone.proto)
	})
}

func TestAddress_ToGraph(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.ToGraph())
	})

	t.Run("strips indices", func(t *testing.T) {
		a := New("aws_instance", "test").CreateIntIndexedChild(0)
		graph := a.ToGraph()
		assert.Equal(t, "aws_instance.test", graph.String())
	})
}

func TestAddress_Truncate(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.Truncate(1))
	})

	t.Run("truncate beyond length", func(t *testing.T) {
		a := New("a", "b")
		result := a.Truncate(10)
		assert.Equal(t, "a.b", result.String())
	})

	t.Run("truncate to 1", func(t *testing.T) {
		a := New("a", "b", "c")
		result := a.Truncate(1)
		assert.Equal(t, "a", result.String())
	})
}

func TestAddress_CreateChild(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		result := a.CreateChild("test")
		assert.Equal(t, "test", result.String())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("parent")
		result := a.CreateChild("child1", "child2")
		assert.Equal(t, "parent.child1.child2", result.String())
	})
}

func TestAddress_CreateIntIndexedChild(t *testing.T) {
	t.Run("empty address", func(t *testing.T) {
		a := New()
		result := a.CreateIntIndexedChild(0)
		assert.Equal(t, "[0]", result.String())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("test")
		result := a.CreateIntIndexedChild(5)
		assert.Equal(t, "test[5]", result.String())
	})

	t.Run("already indexed", func(t *testing.T) {
		a := New("test").CreateIntIndexedChild(0)
		result := a.CreateIntIndexedChild(1)
		assert.Equal(t, "test[0][1]", result.String())
	})
}

func TestAddress_LastIntIndex(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, int64(-1), a.LastIntIndex())
	})

	t.Run("no index", func(t *testing.T) {
		a := New("test")
		assert.Equal(t, int64(-1), a.LastIntIndex())
	})

	t.Run("with index", func(t *testing.T) {
		a := New("test").CreateIntIndexedChild(5)
		assert.Equal(t, int64(5), a.LastIntIndex())
	})
}

func TestAddress_Last(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, "", a.Last())
	})

	t.Run("empty address", func(t *testing.T) {
		a := New()
		assert.Equal(t, "", a.Last())
	})

	t.Run("valid address", func(t *testing.T) {
		a := New("a", "b", "c")
		assert.Equal(t, "c", a.Last())
	})
}

func TestAddress_StartsWith(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		a := New("a", "b")
		assert.True(t, a.StartsWith(New("a", "b")))
	})

	t.Run("prefix match", func(t *testing.T) {
		a := New("a", "b", "c")
		assert.True(t, a.StartsWith(New("a", "b")))
	})

	t.Run("no match", func(t *testing.T) {
		a := New("a", "b")
		assert.False(t, a.StartsWith(New("x", "y")))
	})
}

func TestAddress_StripIndex(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, Empty, a.StripIndex())
	})

	t.Run("with int index", func(t *testing.T) {
		a := New("test").CreateIntIndexedChild(0)
		result := a.StripIndex()
		assert.Equal(t, "test", result.String())
	})

	t.Run("with string index", func(t *testing.T) {
		a := New("test").CreateStringIndexedChild("key")
		result := a.StripIndex()
		assert.Equal(t, "test", result.String())
	})
}

func TestAddress_Hash(t *testing.T) {
	a := New("test", "address")
	hash := a.Hash()
	assert.NotEmpty(t, hash)
	// Same address should produce same hash
	assert.Equal(t, hash, New("test", "address").Hash())
}

func TestAddress_String(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.Equal(t, "", a.String())
	})

	t.Run("empty address", func(t *testing.T) {
		a := New()
		assert.Equal(t, "", a.String())
	})

	t.Run("simple address", func(t *testing.T) {
		a := New("a", "b", "c")
		assert.Equal(t, "a.b.c", a.String())
	})

	t.Run("with int index", func(t *testing.T) {
		a := New("test").CreateIntIndexedChild(0)
		assert.Equal(t, "test[0]", a.String())
	})

	t.Run("with string index", func(t *testing.T) {
		a := New("test").CreateStringIndexedChild("key")
		assert.Equal(t, `test["key"]`, a.String())
	})

	t.Run("standalone int index", func(t *testing.T) {
		a := &Address{proto: &parserpb.Address{
			Segments: []*parserpb.Segment{
				{Value: "a"},
				{IndexInt: int64Ptr(0)},
			},
		}}
		assert.Equal(t, "a[0]", a.String())
	})

	t.Run("standalone string index", func(t *testing.T) {
		a := &Address{proto: &parserpb.Address{
			Segments: []*parserpb.Segment{
				{Value: "a"},
				{IndexString: stringPtr("key")},
			},
		}}
		assert.Equal(t, `a["key"]`, a.String())
	})
}

func TestAddress_RelativeLooksLikeIntIndex(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		var a *Address
		assert.False(t, a.RelativeLooksLikeIntIndex())
	})

	t.Run("empty address", func(t *testing.T) {
		a := New()
		assert.False(t, a.RelativeLooksLikeIntIndex())
	})

	t.Run("starts with int index", func(t *testing.T) {
		a := &Address{proto: &parserpb.Address{
			Segments: []*parserpb.Segment{
				{IndexInt: int64Ptr(0)},
			},
		}}
		assert.True(t, a.RelativeLooksLikeIntIndex())
	})

	t.Run("starts with numeric string", func(t *testing.T) {
		a := &Address{proto: &parserpb.Address{
			Segments: []*parserpb.Segment{
				{Value: "123"},
			},
		}}
		assert.True(t, a.RelativeLooksLikeIntIndex())
	})

	t.Run("starts with non-numeric", func(t *testing.T) {
		a := New("test")
		assert.False(t, a.RelativeLooksLikeIntIndex())
	})
}
