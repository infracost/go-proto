// Package diagnostic provides wrappers and utilities for working with
// diagnostic protocol buffer types for error and warning reporting.
package diagnostic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
)

// Diagnostic represents an error, warning, or informational message from parsing or analysis.
type Diagnostic parserpb.Diagnostic

// ToProto returns the underlying protocol buffer representation.
func (d *Diagnostic) ToProto() *parserpb.Diagnostic {
	if d == nil {
		return nil
	}
	return (*parserpb.Diagnostic)(d)
}

// FormatMessage returns a formatted diagnostic message suitable for display,
// combining the type's message prefix with the error string as
// "<prefix>: <error>".
func (d *Diagnostic) FormatMessage() string {
	return MessagePrefix(d.Type) + ": " + d.Error
}

// NewAsList creates a single diagnostic and returns it as a slice.
func NewAsList(t parserpb.DiagnosticType, format string, params ...any) []*Diagnostic {
	return []*Diagnostic{New(t, format, params...)}
}

// New creates a new diagnostic with the given type and formatted message.
func New(t parserpb.DiagnosticType, format string, params ...any) *Diagnostic {
	return FromError(t, fmt.Errorf(format, params...))
}

// FromError creates a diagnostic from an error with the given type.
func FromError(t parserpb.DiagnosticType, err error) *Diagnostic {
	ignored := isIgnored(err)
	d := &Diagnostic{
		Type:     t,
		Labels:   make(map[string]string),
		Critical: !ignored && IsCritical(t),
		Warning:  !ignored && IsWarning(t),
		Ignored:  ignored,
	}
	if err != nil {
		d.Error = err.Error()
	} else {
		d.Error = "missing error"
	}
	return d
}

var ignoredErrors = []error{
	context.Canceled,
}

func isIgnored(err error) bool {
	if err == nil {
		return false
	}
	for _, ignored := range ignoredErrors {
		if errors.Is(err, ignored) {
			return true
		}
	}
	return false
}

// Surround wraps the diagnostic in a Diagnostics collection.
func (d *Diagnostic) Surround() *Diagnostics {
	return ((*Diagnostics)(nil)).Add(d)
}

// WithSourceRange attaches source location information to the diagnostic.
func (d *Diagnostic) WithSourceRange(sourceRange *parserpb.SourceRange) *Diagnostic {
	d.SourceRange = sourceRange
	return d
}

// WithLabel adds a key-value label to the diagnostic for additional context.
func (d *Diagnostic) WithLabel(key, val string) *Diagnostic {
	d.Labels[key] = val
	return d
}

// String returns a human-readable representation of the diagnostic.
func (d *Diagnostic) String() string {
	if d.SourceRange == nil {
		return fmt.Sprintf("%s: %s %s", d.Type, d.Error, mapToString(d.Labels))
	}
	return fmt.Sprintf("%s: %s %s @ %s", d.Type, d.Error, mapToString(d.Labels), sourceRangeToString(d.SourceRange))
}

func mapToString(m map[string]string) string {
	var s string
	for k, v := range m {
		if s != "" {
			s += ", "
		}
		s += fmt.Sprintf("%s=%v", k, v)
	}
	return s
}

func sourceRangeToString(s *parserpb.SourceRange) string {

	if s == nil {
		return ""
	}

	if s.StartLine == 0 && s.EndLine == 0 && s.StartColumn == 0 && s.EndColumn == 0 && s.Filename == "" {
		return "<unknown>"
	}

	if strings.Contains(s.Filename, "://") {
		return s.Filename
	}

	lineRange := strconv.Itoa(int(s.StartLine))
	if s.StartLine != s.EndLine {
		lineRange += "-" + strconv.Itoa(int(s.EndLine))
	}

	return s.Filename + ":" + lineRange
}
