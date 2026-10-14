package diagnostic

import (
	"testing"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
)

func TestTypesCoverEveryEnumValue(t *testing.T) {
	for v, name := range parserpb.DiagnosticType_name {
		_, ok := typeTable[parserpb.DiagnosticType(v)]
		assert.True(t, ok, "%s has no entry in typeTable", name)
	}
}

func TestTypesNeverBothCriticalAndWarning(t *testing.T) {
	for dt, info := range typeTable {
		assert.False(t, info.critical && info.warning, "%s is both critical and warning", dt)
	}
}

func TestTypeTable(t *testing.T) {
	tests := []struct {
		dt       parserpb.DiagnosticType
		critical bool
		warning  bool
		prefix   string
		code     int
	}{
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_PARSE_ERROR, true, false, "Parse error", CodeJSONParsingFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_EVALUATION_ERROR, true, false, "Evaluation error", CodeModuleEvaluationFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_HELM_CHART, true, false, "Invalid Helm chart", CodeModuleEvaluationFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HELM_RENDER_ERROR, true, false, "Helm render error", CodeModuleEvaluationFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_UNSUPPORTED_CONSTRUCT, false, true, "Unsupported construct", CodeUnsupportedResources},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_INVALID_CLOUDFORMATION_TEMPLATE, true, false, "Invalid CloudFormation template", CodeModuleEvaluationFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MODULE_FETCH_ERROR, true, false, "Failed to load module", CodePrivateModuleDownloadFailure},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_MISSING_INPUT_VARIABLE, false, true, "Missing input variable", CodeMissingVars},
		{parserpb.DiagnosticType_DIAGNOSTIC_TYPE_DEFECT, true, false, "DIAGNOSTIC_TYPE_DEFECT", CodeUnknownFailure},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			assert.Equal(t, tt.critical, IsCritical(tt.dt))
			assert.Equal(t, tt.warning, IsWarning(tt.dt))
			assert.Equal(t, tt.prefix, MessagePrefix(tt.dt))
			assert.Equal(t, tt.code, DashboardCode(tt.dt))
		})
	}
}

func TestDashboardCode_UnknownType(t *testing.T) {
	assert.Equal(t, CodeUnknownFailure, DashboardCode(parserpb.DiagnosticType(9999)))
}
