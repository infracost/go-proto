package diagnostic

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
