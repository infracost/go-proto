// Package synth holds the markers the parser writes into values it cannot know.
package synth

import (
	"regexp"
	"strings"
)

// Identifier stands in for an unset variable.
const Identifier = "infracost-mock-44956be29f34"

// Marker kinds.
const (
	KindPrefix = "prefix"
)

// Marker returns the placeholder for a random value: %ic:<kind>:<id>%.
// It is lowercase so it survives lower().
func Marker(kind, id string) string {
	return "%ic:" + kind + ":" + id + "%"
}

var (
	markerRe     = regexp.MustCompile(`%ic:[a-z0-9_]+:[a-z0-9._-]+%`)
	identifierRe = regexp.MustCompile(`(?i)` + regexp.QuoteMeta(Identifier))
)

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
