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

// messagePrefixes maps each diagnostic type to a human-readable prefix
// used when formatting diagnostic messages for display (e.g. in PR comments).
// The formatted message is: "<prefix>: <error>".
var messagePrefixes = map[parserpb.DiagnosticType]string{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAFORM_CONFIGURATION:      "Invalid Terraform configuration",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAGRUNT_CONFIGURATION:     "Invalid Terragrunt configuration",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR:                   "Failed to load module",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SECURITY_ERROR:                       "Security problem",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP:                   "Invalid source map",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_TERRAFORM_CYCLIC_DEPENDENCY_DETECTED: "Cyclic dependencies detected",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR:                      "HCL parse error",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_FUNC:                         "Missing function",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNKNOWN_EXPRESSION:                   "Unknown HCL expression",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_NO_SOURCE_FILES_FOUND:                "No source files found",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE:               "Missing input variable",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR:                     "Filesystem error",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEPENDENCY_ERROR:                     "Terragrunt dependency error",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_YOR_ERROR:                            "YOR configuration error",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_REMOTE_VARIABLE_LOAD_ERROR:           "Remote variable load error",
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_CLOUDFORMATION_TEMPLATE:      "Invalid CloudFormation template",
}

// MessagePrefix returns the human-readable prefix for the given diagnostic
// type. If the type has no explicit prefix, the proto enum name is returned.
func MessagePrefix(t parserpb.DiagnosticType) string {
	if prefix, ok := messagePrefixes[t]; ok {
		return prefix
	}
	return t.String()
}

// FormatMessage returns a formatted diagnostic message suitable for display,
// combining the type's message prefix with the error string as
// "<prefix>: <error>".
func (d *Diagnostic) FormatMessage() string {
	return MessagePrefix(d.Type) + ": " + d.Error
}

var criticalTypes = map[parserpb.DiagnosticType]struct{}{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSPECIFIED:                          {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT:                               {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR:                   {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAFORM_CONFIGURATION:      {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAGRUNT_CONFIGURATION:     {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_TERRAFORM_CYCLIC_DEPENDENCY_DETECTED: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR:                      {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNKNOWN_EXPRESSION:                   {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SECURITY_ERROR:                       {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_NO_SOURCE_FILES_FOUND:                {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEPENDENCY_ERROR:                     {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP:                   {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_YOR_ERROR:                            {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_CLOUDFORMATION_TEMPLATE:      {},
}

var warningTypes = map[parserpb.DiagnosticType]struct{}{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE:                        {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR:                              {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_REMOTE_VARIABLE_LOAD_ERROR:                    {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CLOUDFORMATION_INTRINSIC_FUNCTION: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CLOUDFORMATION_TRANSFORM:          {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FAILED_CLOUDFORMATION_TRANSFORM:               {},
}

var transientTypes = map[parserpb.DiagnosticType]struct{}{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR:         {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP:         {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_REMOTE_VARIABLE_LOAD_ERROR: {},
}

func isCritical(t parserpb.DiagnosticType) bool {
	_, ok := criticalTypes[t]
	return ok
}

func isWarning(t parserpb.DiagnosticType) bool {
	_, ok := warningTypes[t]
	return ok
}

func isTransient(t parserpb.DiagnosticType) bool {
	_, ok := transientTypes[t]
	return ok
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
		Critical: !ignored && isCritical(t),
		Warning:  !ignored && isWarning(t),
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
