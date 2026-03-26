package diagnostic

import (
	"context"
	"errors"
	"testing"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, "failed to parse %s", "test.tf")
	require.NotNil(t, d)
	assert.Equal(t, parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, d.Type)
	assert.Equal(t, "failed to parse test.tf", d.Error)
	assert.True(t, d.Critical)
}

func TestNewAsList(t *testing.T) {
	list := NewAsList(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "file not found")
	require.Len(t, list, 1)
	assert.Equal(t, "file not found", list[0].Error)
}

func TestFromError(t *testing.T) {
	t.Run("critical error", func(t *testing.T) {
		err := errors.New("parse failed")
		d := FromError(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, err)
		assert.True(t, d.Critical)
		assert.False(t, d.Warning)
		assert.Equal(t, "parse failed", d.Error)
	})

	t.Run("warning error", func(t *testing.T) {
		err := errors.New("missing variable")
		d := FromError(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE, err)
		assert.False(t, d.Critical)
		assert.True(t, d.Warning)
	})

	t.Run("nil error", func(t *testing.T) {
		d := FromError(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, nil)
		assert.Equal(t, "missing error", d.Error)
	})

	t.Run("ignored error", func(t *testing.T) {
		d := FromError(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, context.Canceled)
		assert.True(t, d.Ignored)
		assert.False(t, d.Critical)
	})
}

func TestDiagnostic_ToProto(t *testing.T) {
	t.Run("nil diagnostic", func(t *testing.T) {
		var d *Diagnostic
		assert.Nil(t, d.ToProto())
	})

	t.Run("valid diagnostic", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		proto := d.ToProto()
		require.NotNil(t, proto)
		assert.Equal(t, parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, proto.Type)
	})
}

func TestDiagnostic_Surround(t *testing.T) {
	d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test")
	diags := d.Surround()
	require.NotNil(t, diags)
	assert.Equal(t, 1, diags.Len())
}

func TestDiagnostic_WithSourceRange(t *testing.T) {
	d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test")
	sr := &parserpb.SourceRange{Filename: "test.tf", StartLine: 10, EndLine: 10}
	d = d.WithSourceRange(sr)
	assert.Equal(t, sr, d.SourceRange)
}

func TestDiagnostic_WithLabel(t *testing.T) {
	d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test")
	d = d.WithLabel("module", "vpc")
	assert.Equal(t, "vpc", d.Labels["module"])
}

func TestDiagnostic_String(t *testing.T) {
	t.Run("without source range", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		d.WithLabel("key", "value")
		s := d.String()
		assert.Contains(t, s, "test error")
		assert.Contains(t, s, "DIAGNOSTIC_TYPE_DEFECT")
	})

	t.Run("with source range", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		d.WithSourceRange(&parserpb.SourceRange{Filename: "test.tf", StartLine: 10, EndLine: 12})
		s := d.String()
		assert.Contains(t, s, "test.tf:10-12")
	})

	t.Run("with URL source range", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		d.WithSourceRange(&parserpb.SourceRange{Filename: "https://example.com/test.tf"})
		s := d.String()
		assert.Contains(t, s, "https://example.com/test.tf")
	})

	t.Run("with empty source range", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		d.WithSourceRange(&parserpb.SourceRange{})
		s := d.String()
		assert.Contains(t, s, "<unknown>")
	})

	t.Run("same start and end line", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error")
		d.WithSourceRange(&parserpb.SourceRange{Filename: "test.tf", StartLine: 5, EndLine: 5})
		s := d.String()
		assert.Contains(t, s, "test.tf:5")
		assert.NotContains(t, s, "5-5")
	})
}

func TestMessagePrefix(t *testing.T) {
	t.Run("known type", func(t *testing.T) {
		assert.Equal(t, "HCL parse error", MessagePrefix(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR))
	})

	t.Run("unknown type falls back to enum name", func(t *testing.T) {
		prefix := MessagePrefix(parserpb.DiagnosticType(9999))
		assert.NotEmpty(t, prefix)
	})
}

func TestDiagnostic_FormatMessage(t *testing.T) {
	t.Run("known type", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, "failed to parse test.tf")
		assert.Equal(t, "HCL parse error: failed to parse test.tf", d.FormatMessage())
	})

	t.Run("unknown type uses enum name", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "something broke")
		msg := d.FormatMessage()
		assert.Contains(t, msg, "something broke")
		assert.Contains(t, msg, "DIAGNOSTIC_TYPE_DEFECT")
	})
}

func TestIsCritical(t *testing.T) {
	criticalTypes := []parserpb.DiagnosticType{
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT,
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR,
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SECURITY_ERROR,
	}
	for _, dt := range criticalTypes {
		assert.True(t, isCritical(dt), "expected %v to be critical", dt)
	}
}

func TestIsWarning(t *testing.T) {
	warningTypes := []parserpb.DiagnosticType{
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE,
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR,
	}
	for _, dt := range warningTypes {
		assert.True(t, isWarning(dt), "expected %v to be warning", dt)
	}
}

func TestIsTransient(t *testing.T) {
	transientTypes := []parserpb.DiagnosticType{
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR,
		parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP,
	}
	for _, dt := range transientTypes {
		assert.True(t, isTransient(dt), "expected %v to be transient", dt)
	}
}
