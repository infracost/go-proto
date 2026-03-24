package resource

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func val(s string) value.String {
	return value.New(s, 0, "", nil)
}

func TestGetBase(t *testing.T) {
	r := &Resource{ID: "test-id"}
	assert.Equal(t, r, r.GetBase())
}

func TestTags_Get(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("prod")},
		{Key: val("team"), Value: val("platform")},
	}

	v, ok := tags.Get("env")
	require.True(t, ok)
	assert.Equal(t, "prod", v.Value())

	v, ok = tags.Get("team")
	require.True(t, ok)
	assert.Equal(t, "platform", v.Value())

	_, ok = tags.Get("missing")
	assert.False(t, ok)
}

func TestTags_Set_NewTag(t *testing.T) {
	var tags Tags
	tags.Set(val("env"), val("prod"), false)

	require.Len(t, tags, 1)
	assert.Equal(t, "env", tags[0].Key.Value())
	assert.Equal(t, "prod", tags[0].Value.Value())
	assert.False(t, tags[0].IsDefault)
}

func TestTags_Set_OverwriteDefaultWithNonDefault(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("default-val"), IsDefault: true},
	}

	tags.Set(val("env"), val("explicit-val"), false)

	require.Len(t, tags, 1)
	assert.Equal(t, "explicit-val", tags[0].Value.Value())
	assert.False(t, tags[0].IsDefault)
}

func TestTags_Set_NonDefaultNotOverwrittenByDefault(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("explicit-val"), IsDefault: false},
	}

	tags.Set(val("env"), val("default-val"), true)

	require.Len(t, tags, 1)
	assert.Equal(t, "explicit-val", tags[0].Value.Value())
	assert.False(t, tags[0].IsDefault)
}

func TestTags_Set_NonDefaultOverwrittenByNonDefault(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("old"), IsDefault: false},
	}

	tags.Set(val("env"), val("new"), false)

	require.Len(t, tags, 1)
	assert.Equal(t, "new", tags[0].Value.Value())
}

func TestTags_DefaultChecksum_OnlyIncludesDefaults(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("prod"), IsDefault: true},
		{Key: val("team"), Value: val("platform"), IsDefault: false},
		{Key: val("cost-center"), Value: val("123"), IsDefault: true},
	}

	checksum1 := tags.DefaultChecksum()
	assert.NotEmpty(t, checksum1)

	// same default tags should produce same checksum
	tags2 := Tags{
		{Key: val("cost-center"), Value: val("123"), IsDefault: true},
		{Key: val("env"), Value: val("prod"), IsDefault: true},
		{Key: val("team"), Value: val("other"), IsDefault: false},
	}

	checksum2 := tags2.DefaultChecksum()
	assert.Equal(t, checksum1, checksum2)
}

func TestTags_DefaultChecksum_NoDefaults(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("prod"), IsDefault: false},
	}

	// checksum with no defaults should equal checksum of empty string
	checksum := tags.DefaultChecksum()
	emptyChecksum := Tags{}.DefaultChecksum()
	assert.Equal(t, emptyChecksum, checksum)
}

func TestTags_DefaultChecksum_DifferentValues(t *testing.T) {
	tags1 := Tags{
		{Key: val("env"), Value: val("prod"), IsDefault: true},
	}
	tags2 := Tags{
		{Key: val("env"), Value: val("staging"), IsDefault: true},
	}

	assert.NotEqual(t, tags1.DefaultChecksum(), tags2.DefaultChecksum())
}

func TestTags_ToProto(t *testing.T) {
	tags := Tags{
		{Key: val("env"), Value: val("prod"), IsDefault: true},
		{Key: val("team"), Value: val("platform"), IsDefault: false},
	}

	proto := tags.ToProto()
	require.Len(t, proto, 2)
	assert.Equal(t, "env", proto[0].Key.GetStringValue())
	assert.Equal(t, "prod", proto[0].Value.GetStringValue())
	assert.True(t, proto[0].IsDefault)
	assert.Equal(t, "team", proto[1].Key.GetStringValue())
	assert.Equal(t, "platform", proto[1].Value.GetStringValue())
	assert.False(t, proto[1].IsDefault)
}

func TestTagsFromProto(t *testing.T) {
	original := Tags{
		{Key: val("env"), Value: val("prod"), IsDefault: true},
		{Key: val("team"), Value: val("platform"), IsDefault: false},
	}

	proto := original.ToProto()
	roundTripped := TagsFromProto(proto)

	require.Len(t, roundTripped, 2)
	assert.Equal(t, "env", roundTripped[0].Key.Value())
	assert.Equal(t, "prod", roundTripped[0].Value.Value())
	assert.True(t, roundTripped[0].IsDefault)
	assert.Equal(t, "team", roundTripped[1].Key.Value())
	assert.Equal(t, "platform", roundTripped[1].Value.Value())
	assert.False(t, roundTripped[1].IsDefault)
}
