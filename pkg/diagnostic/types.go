package diagnostic

import (
	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
)

// typeInfo holds everything known about one diagnostic type. A type that is
// neither critical nor warning is informational, an empty prefix shows the
// enum name instead, and a zero code is CodeUnknownFailure.
type typeInfo struct {
	critical  bool
	warning   bool
	transient bool
	prefix    string
	code      int
}

// typeTable has an entry for every DiagnosticType, so adding one to the proto
// fails TestTypesCoverEveryEnumValue until it is classified here.
var typeTable = map[parserpb.DiagnosticType]typeInfo{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSPECIFIED:      {critical: true},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FAILED_OPERATION: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT:           {critical: true},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR: {
		critical: true, prefix: "HCL parse error", code: CodeHCLParseError,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_FUNC: {
		prefix: "Missing function", code: CodeMissingFunc,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FAILED_FUNCTION_CALL: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNEXPECTED_TYPE:      {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR: {
		critical: true, transient: true, prefix: "Failed to load module", code: CodePrivateModuleDownloadFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE: {
		warning: true, prefix: "Missing input variable", code: CodeMissingVars,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_REFERENCE: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_VALUE_INDEX_ERROR: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNKNOWN_EXPRESSION: {
		critical: true, prefix: "Unknown HCL expression", code: CodeUnknownExpression,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_RELATIVE_TRAVERSAL_ERROR: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAFORM_CONFIGURATION: {
		critical: true, prefix: "Invalid Terraform configuration", code: CodeModuleEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAGRUNT_CONFIGURATION: {
		critical: true, prefix: "Invalid Terragrunt configuration", code: CodeTerragruntEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_TERRAFORM_CYCLIC_DEPENDENCY_DETECTED: {
		critical: true, prefix: "Cyclic dependencies detected", code: CodeCyclicDependenciesDetected,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_TERRAGRUNT_CYCLIC_DEPENDENCY_DETECTED: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SECURITY_ERROR: {
		critical: true, prefix: "Security problem", code: CodeSecurityError,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEPENDENCY_ERROR: {
		critical: true, prefix: "Terragrunt dependency error", code: CodeTerragruntEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR: {
		warning: true, prefix: "Filesystem error", code: CodeFilesystemError,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP: {
		critical: true, transient: true, prefix: "Invalid source map", code: CodeInvalidSourceMap,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_NO_SOURCE_FILES_FOUND: {
		critical: true, prefix: "No source files found", code: CodeEmptyPath,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_YOR_ERROR: {
		critical: true, prefix: "YOR configuration error", code: CodeYORError,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_REMOTE_VARIABLE_LOAD_ERROR: {
		warning: true, transient: true, prefix: "Remote variable load error", code: CodeRemoteVariableLoadError,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_CLOUDFORMATION_TEMPLATE: {
		critical: true, prefix: "Invalid CloudFormation template", code: CodeModuleEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CLOUDFORMATION_TRANSFORM: {
		warning: true, prefix: "Unsupported CloudFormation transform",
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CLOUDFORMATION_INTRINSIC_FUNCTION: {
		warning: true, prefix: "Unsupported CloudFormation intrinsic function",
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SYNTHETIC_CLOUDFORMATION_TEMPLATE_URL: {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FAILED_CLOUDFORMATION_TRANSFORM: {
		warning: true, prefix: "Failed CloudFormation transform",
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_LOADER_OPTION_ERROR:           {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INPUT_TARGET_ERROR:            {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INPUT_TARGET_OPTIONS_ERROR:    {},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_BUILDING_LOADER_OPTIONS_ERROR: {},
	// The cloud-agnostic forms of the per-IaC parse and evaluation types. A
	// source that fails to parse, or a template that fails structural
	// validation, would produce a result built from garbage.
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_PARSE_ERROR: {
		critical: true, prefix: "Parse error", code: CodeJSONParsingFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_EVALUATION_ERROR: {
		critical: true, prefix: "Evaluation error", code: CodeModuleEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_CYCLIC_DEPENDENCY_DETECTED: {},
	// Recognized but not modelled: the rest still parses and costs, so this
	// warns that resources may be missing rather than failing the run.
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CONSTRUCT: {
		warning: true, prefix: "Unsupported construct", code: CodeUnsupportedResources,
	},
	// A chart Helm refuses to load or render produces no workloads at all.
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_HELM_CHART: {
		critical: true, prefix: "Invalid Helm chart", code: CodeModuleEvaluationFailure,
	},
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HELM_RENDER_ERROR: {
		critical: true, prefix: "Helm render error", code: CodeModuleEvaluationFailure,
	},
}

// IsCritical reports whether diagnostics of type t fail the run.
func IsCritical(t parserpb.DiagnosticType) bool {
	return typeTable[t].critical
}

// IsWarning reports whether diagnostics of type t are shown as warnings.
func IsWarning(t parserpb.DiagnosticType) bool {
	return typeTable[t].warning
}

// IsTransient reports whether diagnostics of type t may resolve on retry.
func IsTransient(t parserpb.DiagnosticType) bool {
	return typeTable[t].transient
}

// MessagePrefix returns the human-readable prefix for the given diagnostic
// type. If the type has no explicit prefix, the proto enum name is returned.
func MessagePrefix(t parserpb.DiagnosticType) string {
	if prefix := typeTable[t].prefix; prefix != "" {
		return prefix
	}
	return t.String()
}

// DashboardCode returns the dashboard-specific integer code for the given
// diagnostic type. If the type has no explicit mapping, CodeUnknownFailure
// is returned.
func DashboardCode(t parserpb.DiagnosticType) int {
	return typeTable[t].code
}
