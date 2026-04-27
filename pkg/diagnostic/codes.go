package diagnostic

import (
	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
)

// Dashboard diagnostic codes — these are integer codes used by the Infracost
// dashboard to categorize and display project-level diagnostics. They are
// stable and must not be renumbered.
const (
	// Local module issues
	CodeUnknownFailure                    = 0
	CodeJSONParsingFailure                = 101
	CodeModuleEvaluationFailure           = 102
	CodeTerragruntEvaluationFailure       = 103
	CodeTerragruntModuleEvaluationFailure = 104
	CodeMissingVars                       = 105
	CodeEmptyPath                         = 106
	CodeInvalidSourceMap                  = 107
	CodeCyclicDependenciesDetected        = 108
	CodeHCLParseError                     = 109
	CodeMissingFunc                       = 110
	CodeFilesystemError                   = 111
	CodeUnknownExpression                 = 112

	// Git module issues
	CodePrivateModuleDownloadFailure = 201

	// Infracost cloud issues
	CodeRunQuotaExceeded = 401
	CodeInvalidAPIKey    = 402

	// Security / policy issues
	CodeSecurityError          = 501
	CodePolicyEvaluationFailed = 502

	// CDK errors
	CodeCDKSynthFailure    = 701
	CodeCDKRegionDefaulted = 702

	// Other
	CodeBuildScriptError        = 800
	CodeUnsupportedResources    = 801
	CodeRemoteVariableLoadError = 802
	CodeYORError                = 803
)

// DashboardCode returns the dashboard-specific integer code for the given
// diagnostic type. If the type has no explicit mapping, CodeUnknownFailure
// is returned.
func DashboardCode(t parserpb.DiagnosticType) int {
	if code, ok := dashboardCodeMap[t]; ok {
		return code
	}
	return CodeUnknownFailure
}

var dashboardCodeMap = map[parserpb.DiagnosticType]int{
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAFORM_CONFIGURATION:      CodeModuleEvaluationFailure,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_TERRAGRUNT_CONFIGURATION:     CodeTerragruntEvaluationFailure,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR:                   CodePrivateModuleDownloadFailure,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_SECURITY_ERROR:                       CodeSecurityError,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_SOURCE_MAP:                   CodeInvalidSourceMap,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_TERRAFORM_CYCLIC_DEPENDENCY_DETECTED: CodeCyclicDependenciesDetected,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR:                      CodeHCLParseError,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_FUNC:                         CodeMissingFunc,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNKNOWN_EXPRESSION:                   CodeUnknownExpression,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_NO_SOURCE_FILES_FOUND:                CodeEmptyPath,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE:               CodeMissingVars,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR:                     CodeFilesystemError,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEPENDENCY_ERROR:                     CodeTerragruntEvaluationFailure,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_YOR_ERROR:                            CodeYORError,
	parserpb.DiagnosticType_DIAGNOSTIC_TYPE_REMOTE_VARIABLE_LOAD_ERROR:           CodeRemoteVariableLoadError,
}