package event

import (
	"strings"

	"github.com/infracost/proto/gen/go/infracost/parser/event"
)

type StringFilter struct {
	*event.StringFilter
}

func StringFilterFromProto(pb *event.StringFilter) *StringFilter {
	return &StringFilter{pb}
}

func (f *StringFilter) Match(projectName string) bool {
	if len(f.GetExclude()) > 0 {
		for _, exclude := range f.GetExclude() {
			if matchWildcard(projectName, exclude) {
				return false
			}
		}
	}
	if len(f.GetInclude()) > 0 {
		for _, include := range f.GetInclude() {
			if matchWildcard(projectName, include) {
				return true
			}
		}
		return false
	}
	return true
}

func (f *StringFilter) Matches(value string) bool {

	if f == nil || f.StringFilter == nil {
		return true
	}

	// 1. if an exclude rule matches, return false
	// 2. if there are any include rules, return whether one matches or not
	// 3. if there are no include rules, return true

	for _, exclude := range f.GetExclude() {
		if matchWildcard(value, exclude) {
			return false
		}
	}

	if len(f.GetInclude()) > 0 {
		for _, include := range f.GetInclude() {
			if matchWildcard(value, include) {
				return true
			}
		}
		return false
	}

	return true
}

// MatchesText applies the filter to a block of text, line by line. An entry
// matches a line that contains it, so "CET-" matches "Approved under CET-1234",
// and wildcards still work within a line. Matching is case-insensitive.
func (f *StringFilter) MatchesText(text string) bool {
	if f == nil || f.StringFilter == nil {
		return true
	}

	if len(f.GetExclude()) == 0 && len(f.GetInclude()) == 0 {
		return true
	}

	var lines []string
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, strings.ToLower(trimmed))
		}
	}

	for _, exclude := range f.GetExclude() {
		if anyLineContains(lines, exclude) {
			return false
		}
	}

	if len(f.GetInclude()) > 0 {
		for _, include := range f.GetInclude() {
			if anyLineContains(lines, include) {
				return true
			}
		}
		return false
	}

	return true
}

// anyLineContains reports whether entry appears anywhere in any line. Lines are
// already lowercased. An empty entry matches nothing, as it does in Matches.
func anyLineContains(lines []string, entry string) bool {
	entry = strings.ToLower(strings.TrimSpace(entry))
	if entry == "" {
		return false
	}

	pattern := "*" + entry + "*"
	for _, line := range lines {
		if matchWildcard(line, pattern) {
			return true
		}
	}
	return false
}

// matchWildcard checks if a value matches a pattern with wildcards
// The pattern can contain '*' as a wildcard to match any number of characters
func matchWildcard(value string, pattern string) bool {
	if pattern == "*" {
		return true
	}

	patternLen := len(pattern)
	valueLen := len(value)
	var patternIndex, valueIndex, starIndex, matchIndex int
	starIndex = -1
	matchIndex = -1

	for valueIndex < valueLen {
		switch {
		case patternIndex < patternLen && pattern[patternIndex] == '*':
			// Found a wildcard, mark position
			starIndex = patternIndex
			matchIndex = valueIndex
			patternIndex++
		case patternIndex < patternLen && (pattern[patternIndex] == value[valueIndex] || pattern[patternIndex] == '?'):
			// Characters match or single character wildcard
			patternIndex++
			valueIndex++
		case starIndex != -1:
			// No match, but we have a previous wildcard
			patternIndex = starIndex + 1
			matchIndex++
			valueIndex = matchIndex
		default:
			// No match and no wildcard
			return false
		}
	}

	// Handle any remaining pattern characters
	for patternIndex < patternLen && pattern[patternIndex] == '*' {
		patternIndex++
	}

	return patternIndex == patternLen
}
