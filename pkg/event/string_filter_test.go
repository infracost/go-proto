package event

import (
	"fmt"
	"testing"

	eventpb "github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/stretchr/testify/assert"
)

func Test_matchWildcard(t *testing.T) {
	tests := []struct {
		value    string
		pattern  string
		expected bool
	}{
		{
			value:    "test",
			pattern:  "test",
			expected: true,
		},
		{
			value:    "test",
			pattern:  "test*",
			expected: true,
		},
		{
			value:    "test",
			pattern:  "*test",
			expected: true,
		},
		{
			value:    "test",
			pattern:  "t*st",
			expected: true,
		},
		{
			value:    "test",
			pattern:  "t*st",
			expected: true,
		},
		{
			value:    "test",
			pattern:  "te*st",
			expected: true,
		},
		{
			value:    "txst",
			pattern:  "tes*",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("matchWildcard(%s %s)", test.value, test.pattern), func(t *testing.T) {
			result := matchWildcard(test.value, test.pattern)
			assert.Equal(t, test.expected, result)
		})
	}

}

func TestStringFilter_Matches(t *testing.T) {
	tests := []struct {
		name   string
		filter *StringFilter
		value  string
		want   bool
	}{
		{
			name:   "nil filter matches everything",
			filter: nil,
			value:  "anything",
			want:   true,
		},
		{
			name:   "nil proto matches everything",
			filter: &StringFilter{nil},
			value:  "anything",
			want:   true,
		},
		{
			name:   "empty filter matches everything",
			filter: StringFilterFromProto(&eventpb.StringFilter{}),
			value:  "anything",
			want:   true,
		},
		{
			name:   "include matches exact value",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_instance"}}),
			value:  "aws_instance",
			want:   true,
		},
		{
			name:   "include rejects non-matching value",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_instance"}}),
			value:  "aws_s3_bucket",
			want:   false,
		},
		{
			name:   "include with wildcard matches prefix",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_*"}}),
			value:  "aws_instance",
			want:   true,
		},
		{
			name:   "include with wildcard rejects non-matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_*"}}),
			value:  "google_compute_instance",
			want:   false,
		},
		{
			name:   "multiple includes - first matches",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_instance", "aws_s3_bucket"}}),
			value:  "aws_instance",
			want:   true,
		},
		{
			name:   "multiple includes - second matches",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_instance", "aws_s3_bucket"}}),
			value:  "aws_s3_bucket",
			want:   true,
		},
		{
			name:   "multiple includes - none match",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"aws_instance", "aws_s3_bucket"}}),
			value:  "aws_lambda_function",
			want:   false,
		},
		{
			name:   "exclude rejects matching value",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"aws_instance"}}),
			value:  "aws_instance",
			want:   false,
		},
		{
			name:   "exclude allows non-matching value",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"aws_instance"}}),
			value:  "aws_s3_bucket",
			want:   true,
		},
		{
			name:   "exclude with wildcard rejects matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"aws_*"}}),
			value:  "aws_instance",
			want:   false,
		},
		{
			name:   "exclude with wildcard allows non-matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"aws_*"}}),
			value:  "google_compute_instance",
			want:   true,
		},
		{
			name: "include and exclude - exclude takes precedence",
			filter: StringFilterFromProto(&eventpb.StringFilter{
				Include: []string{"aws_*"},
				Exclude: []string{"aws_instance"},
			}),
			value: "aws_instance",
			want:  false,
		},
		{
			name: "include and exclude - not excluded, included",
			filter: StringFilterFromProto(&eventpb.StringFilter{
				Include: []string{"aws_*"},
				Exclude: []string{"aws_instance"},
			}),
			value: "aws_s3_bucket",
			want:  true,
		},
		{
			name: "include and exclude - not excluded, not included",
			filter: StringFilterFromProto(&eventpb.StringFilter{
				Include: []string{"aws_*"},
				Exclude: []string{"aws_instance"},
			}),
			value: "google_compute_instance",
			want:  false,
		},
		{
			name:   "empty string matches wildcard include",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"*"}}),
			value:  "",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.filter.Matches(tt.value))
		})
	}
}

func TestStringFilter_Match(t *testing.T) {
	tests := []struct {
		name   string
		filter *StringFilter
		value  string
		want   bool
	}{
		{
			name:   "empty filter matches everything",
			filter: StringFilterFromProto(&eventpb.StringFilter{}),
			value:  "my-project",
			want:   true,
		},
		{
			name:   "include matches exact",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"my-project"}}),
			value:  "my-project",
			want:   true,
		},
		{
			name:   "include rejects non-matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"my-project"}}),
			value:  "other-project",
			want:   false,
		},
		{
			name:   "include with wildcard",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"infra-*"}}),
			value:  "infra-core",
			want:   true,
		},
		{
			name:   "exclude rejects matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"test-*"}}),
			value:  "test-project",
			want:   false,
		},
		{
			name:   "exclude allows non-matching",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"test-*"}}),
			value:  "prod-project",
			want:   true,
		},
		{
			name: "exclude takes precedence over include",
			filter: StringFilterFromProto(&eventpb.StringFilter{
				Include: []string{"*"},
				Exclude: []string{"secret-*"},
			}),
			value: "secret-project",
			want:  false,
		},
		{
			name: "not excluded and included passes",
			filter: StringFilterFromProto(&eventpb.StringFilter{
				Include: []string{"*"},
				Exclude: []string{"secret-*"},
			}),
			value: "public-project",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.filter.Match(tt.value))
		})
	}
}

func TestStringFilterFromProto(t *testing.T) {
	t.Run("wraps proto correctly", func(t *testing.T) {
		pb := &eventpb.StringFilter{
			Include: []string{"a", "b"},
			Exclude: []string{"c"},
		}
		f := StringFilterFromProto(pb)
		assert.Equal(t, []string{"a", "b"}, f.GetInclude())
		assert.Equal(t, []string{"c"}, f.GetExclude())
	})

	t.Run("nil proto is handled", func(t *testing.T) {
		f := StringFilterFromProto(nil)
		assert.True(t, f.Matches("anything"))
	})
}

func TestStringFilter_MatchesText(t *testing.T) {
	const body = "Bumps the cluster size.\n\nApproved under CET-1234, see the ticket.\n"

	tests := []struct {
		name   string
		filter *StringFilter
		text   string
		want   bool
	}{
		{
			name:   "nil filter matches",
			filter: StringFilterFromProto(nil),
			text:   body,
			want:   true,
		},
		{
			name:   "empty filter matches",
			filter: StringFilterFromProto(&eventpb.StringFilter{}),
			text:   body,
			want:   true,
		},
		{
			name:   "exclude entry found mid-line does not match",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-"}}),
			text:   body,
			want:   false,
		},
		{
			name:   "exclude entry absent matches",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-"}}),
			text:   "Bumps the cluster size.",
			want:   true,
		},
		{
			name:   "exclude is case-insensitive",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-"}}),
			text:   "approved under cet-1234",
			want:   false,
		},
		{
			name:   "wildcards work within a line",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-*-approved"}}),
			text:   "ticket CET-1234-approved by finance",
			want:   false,
		},
		{
			name:   "include entry found matches",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"CET-"}}),
			text:   body,
			want:   true,
		},
		{
			name:   "include entry absent does not match",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"CET-"}}),
			text:   "Bumps the cluster size.",
			want:   false,
		},
		{
			name:   "exclude wins over include",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"cluster"}, Exclude: []string{"CET-"}}),
			text:   body,
			want:   false,
		},
		{
			name:   "empty text cannot satisfy an include",
			filter: StringFilterFromProto(&eventpb.StringFilter{Include: []string{"CET-"}}),
			text:   "",
			want:   false,
		},
		{
			name:   "empty text cannot trip an exclude",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-"}}),
			text:   "",
			want:   true,
		},
		{
			name:   "a match on one line is enough",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"CET-"}}),
			text:   "line one\nline two\nCET-9\nline four",
			want:   false,
		},
		{
			name:   "entries are trimmed",
			filter: StringFilterFromProto(&eventpb.StringFilter{Exclude: []string{"  CET-  "}}),
			text:   body,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.filter.MatchesText(tt.text))
		})
	}
}
