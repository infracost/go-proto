package diagnostic

// ErrorCode is used by the dashboard to handle certain run-level errors.
// These codes are referenced in the dashboard, so be careful about changing values.
type ErrorCode uint16

// Run-level error codes.
const (
	ErrorCodeCLIError                         ErrorCode = 300
	ErrorCodeVCSCloneError                    ErrorCode = 301
	ErrorCodeVCSFetchError                    ErrorCode = 302
	ErrorCodeCLIUploadError                   ErrorCode = 303
	ErrorCodeCLICommentError                  ErrorCode = 304
	ErrorCodeCLIOutputError                   ErrorCode = 305
	ErrorCodeCLIGenerateConfigError           ErrorCode = 306
	ErrorCodeCLIBreakdownError                ErrorCode = 307
	ErrorCodeCLIDiffError                     ErrorCode = 308
	ErrorCodeInvalidConfig                    ErrorCode = 309
	ErrorCodeInvalidConfigTemplate            ErrorCode = 310
	ErrorCodeCLIPanicError                    ErrorCode = 320
	ErrorCodeCDKError                         ErrorCode = 330
	ErrorCodeVCSRefMissing                    ErrorCode = 350
	ErrorCodeVCSRefInvalid                    ErrorCode = 351
	ErrorCodeVCSEmptyRepo                     ErrorCode = 352
	ErrorCodeVCSError                         ErrorCode = 360
	ErrorCodeVCSAuthError                     ErrorCode = 365
	ErrorCodeGraphQLError                     ErrorCode = 370
	ErrorCodeGraphQLDisconnectionError        ErrorCode = 371
	ErrorCodeBuildScriptError                 ErrorCode = 800
	ErrorCodeCustomPropertyMappingScriptError ErrorCode = 801
	ErrorCodeCommandKilled                    ErrorCode = 996
	ErrorCodeNoIaC                            ErrorCode = 997
	ErrorCodeInternalError                    ErrorCode = 998
	ErrorCodeUnknownError                     ErrorCode = 999
)