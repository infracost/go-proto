// Package synth holds the markers the parser writes into values it cannot know.
package synth

import (
	"regexp"
	"strconv"
	"strings"
)

// Identifier stands in for an unset variable.
const Identifier = "infracost-mock-44956be29f34"

// Marker kinds.
const (
	KindPrefix = "prefix"
)

// DefaultLength is the length assumed for a marker that does not carry one.
const DefaultLength = 8

// Marker returns the placeholder for a random value: %ic:<kind>:<id>%.
// It is lowercase so it survives lower().
func Marker(kind, id string) string {
	return "%ic:" + kind + ":" + id + "%"
}

// MarkerLen is Marker for a value whose length is known: %ic:<kind>:<id>:<n>%.
func MarkerLen(kind, id string, n int) string {
	return "%ic:" + kind + ":" + id + ":" + strconv.Itoa(n) + "%"
}

var (
	markerRe     = regexp.MustCompile(`(?i)%ic:[a-z0-9_]+:[a-z0-9._-]+?(?::([0-9]+))?%`)
	identifierRe = regexp.MustCompile(`(?i)` + regexp.QuoteMeta(Identifier))
)

// Found is one marker in a string.
type Found struct {
	Start, End int // byte offsets of the marker in the string
	Length     int // the length of the value it stands for
}

// Markers returns the markers in s, in order.
func Markers(s string) []Found {
	var out []Found
	for _, m := range markerRe.FindAllStringSubmatchIndex(s, -1) {
		n := DefaultLength
		if m[2] >= 0 {
			if v, err := strconv.Atoi(s[m[2]:m[3]]); err == nil {
				n = v
			}
		}
		out = append(out, Found{Start: m[0], End: m[1], Length: n})
	}
	return out
}

// minLiteral is the fewest literal chars a pattern must keep to be useful.
const minLiteral = 4

// Wildcard replaces every marker and Identifier in s with '*'.
func Wildcard(s string) string {
	s = markerRe.ReplaceAllString(s, "*")
	s = identifierRe.ReplaceAllString(s, "*")
	for strings.Contains(s, "**") {
		s = strings.ReplaceAll(s, "**", "*")
	}
	return s
}

// Pattern returns s as a '*' glob, or "" if s holds no unknown part or keeps
// fewer than four literal chars.
func Pattern(s string) string {
	p := Wildcard(s)
	if !strings.Contains(p, "*") {
		return ""
	}
	if len(strings.ReplaceAll(p, "*", "")) < minLiteral {
		return ""
	}
	return p
}
