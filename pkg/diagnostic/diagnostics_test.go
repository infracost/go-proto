package diagnostic

import (
	"testing"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiagnostics(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		diags := NewDiagnostics(nil)
		require.NotNil(t, diags)
		assert.Equal(t, 0, diags.Len())
	})

	t.Run("with diagnostics", func(t *testing.T) {
		d1 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "error 1")
		d2 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "error 2")
		diags := NewDiagnostics([]*Diagnostic{d1, d2})
		assert.Equal(t, 2, diags.Len())
	})
}

func TestFromProto(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		diags := FromProto(nil)
		require.NotNil(t, diags)
		assert.Equal(t, 0, diags.Len())
	})

	t.Run("with protos", func(t *testing.T) {
		protos := []*parserpb.Diagnostic{
			{Type: parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, Error: "test"},
		}
		diags := FromProto(protos)
		assert.Equal(t, 1, diags.Len())
	})
}

func TestDiagnostics_Add(t *testing.T) {
	t.Run("add to nil", func(t *testing.T) {
		var diags *Diagnostics
		diags = diags.Add(New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test"))
		require.NotNil(t, diags)
		assert.Equal(t, 1, diags.Len())
	})

	t.Run("add critical", func(t *testing.T) {
		diags := &Diagnostics{}
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, "critical")
		diags = diags.Add(d)
		assert.Equal(t, 1, diags.Len())
	})

	t.Run("limit non-critical", func(t *testing.T) {
		diags := &Diagnostics{}
		// Add more than maxNonCriticalDiagnostics non-critical items
		for i := 0; i < 15; i++ {
			d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "warning")
			d.Critical = false
			diags = diags.Add(d)
		}
		// Should be capped at maxNonCriticalDiagnostics (10)
		assert.Equal(t, 10, diags.Len())
	})

	t.Run("critical not limited", func(t *testing.T) {
		diags := &Diagnostics{}
		for i := 0; i < 15; i++ {
			d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "critical")
			diags = diags.Add(d)
		}
		assert.Equal(t, 15, diags.Len())
	})
}

func TestDiagnostics_Merge(t *testing.T) {
	t.Run("merge into nil", func(t *testing.T) {
		var diags *Diagnostics
		other := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test"),
		})
		result := diags.Merge(other)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Len())
	})

	t.Run("merge nil", func(t *testing.T) {
		diags := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test"),
		})
		result := diags.Merge(nil)
		assert.Equal(t, 1, result.Len())
	})

	t.Run("merge two", func(t *testing.T) {
		diags1 := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test1"),
		})
		diags2 := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test2"),
		})
		result := diags1.Merge(diags2)
		assert.Equal(t, 2, result.Len())
	})
}

func TestDiagnostics_Unwrap(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.Unwrap())
	})

	t.Run("with diagnostics", func(t *testing.T) {
		d := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test")
		diags := NewDiagnostics([]*Diagnostic{d})
		list := diags.Unwrap()
		require.Len(t, list, 1)
		assert.Equal(t, d, list[0])
	})
}

func TestDiagnostics_ToProto(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.ToProto())
	})

	t.Run("with diagnostics", func(t *testing.T) {
		diags := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test"),
		})
		protos := diags.ToProto()
		require.Len(t, protos, 1)
		assert.Equal(t, "test", protos[0].Error)
	})
}

func TestDiagnostics_String(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Equal(t, "", diags.String())
	})

	t.Run("with diagnostics", func(t *testing.T) {
		diags := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test error"),
		})
		s := diags.String()
		assert.Contains(t, s, "test error")
	})
}

func TestDiagnostics_OfType(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.OfType(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT))
	})

	t.Run("filter by type", func(t *testing.T) {
		diags := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "defect"),
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "fs error"),
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "defect 2"),
		})
		filtered := diags.OfType(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT)
		require.NotNil(t, filtered)
		assert.Equal(t, 2, filtered.Len())
	})
}

func TestDiagnostics_Critical(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.Critical())
	})

	t.Run("filter critical", func(t *testing.T) {
		d1 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "critical")
		d2 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "warning")
		d2.Critical = false
		diags := NewDiagnostics([]*Diagnostic{d1, d2})
		critical := diags.Critical()
		require.NotNil(t, critical)
		assert.Equal(t, 1, critical.Len())
	})
}

func TestDiagnostics_NonCritical(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.NonCritical())
	})

	t.Run("filter non-critical", func(t *testing.T) {
		d1 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "critical")
		d2 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, "warning")
		d2.Critical = false
		diags := NewDiagnostics([]*Diagnostic{d1, d2})
		nonCritical := diags.NonCritical()
		require.NotNil(t, nonCritical)
		assert.Equal(t, 1, nonCritical.Len())
	})
}

func TestDiagnostics_Ignored(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.Ignored())
	})

	t.Run("filter ignored", func(t *testing.T) {
		d1 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "not ignored")
		d2 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "ignored")
		d2.Ignored = true
		d2.Critical = false
		diags := NewDiagnostics([]*Diagnostic{d1, d2})
		ignored := diags.Ignored()
		require.NotNil(t, ignored)
		assert.Equal(t, 1, ignored.Len())
	})
}

func TestDiagnostics_Transient(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Nil(t, diags.Transient())
	})

	t.Run("filter transient", func(t *testing.T) {
		d1 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "not transient")
		d2 := New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR, "transient")
		diags := NewDiagnostics([]*Diagnostic{d1, d2})
		transient := diags.Transient()
		require.NotNil(t, transient)
		assert.Equal(t, 1, transient.Len())
	})
}

func TestDiagnostics_Len(t *testing.T) {
	t.Run("nil diagnostics", func(t *testing.T) {
		var diags *Diagnostics
		assert.Equal(t, 0, diags.Len())
	})

	t.Run("with diagnostics", func(t *testing.T) {
		diags := NewDiagnostics([]*Diagnostic{
			New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, "test"),
		})
		assert.Equal(t, 1, diags.Len())
	})
}
