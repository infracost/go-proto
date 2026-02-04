// Package address provides utilities for working with Terraform addresses.
// An address represents a path to a resource, variable, or other entity in
// Terraform configuration, composed of segments that can include names and indices.
package address

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"

	parserpb "github.com/infracost/proto/gen/go/infracost/parser"
)

// Empty is a pre-allocated empty address for convenience.
var Empty = &Address{proto: &parserpb.Address{}}

// Address represents a Terraform address, wrapping the protocol buffer type
// with utility methods for manipulation and comparison.
type Address struct {
	proto *parserpb.Address
	str   *string
	strMu sync.RWMutex
}

// ToProto returns the underlying protocol buffer representation.
func (a *Address) ToProto() *parserpb.Address {
	if a == nil || a.proto == nil {
		return nil
	}
	return a.proto
}

// FromProto creates an Address from a protocol buffer representation.
func FromProto(proto *parserpb.Address) *Address {
	if proto == nil {
		return Empty
	}
	return &Address{proto: proto}
}

// MarshalJSON implements json.Marshaler.
func (a *Address) MarshalJSON() ([]byte, error) {
	if a == nil || a.proto == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(a.proto)
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Address) UnmarshalJSON(data []byte) error {
	var proto parserpb.Address
	if err := json.Unmarshal(data, &proto); err != nil {
		return err
	}
	a.proto = &proto
	return nil
}

// Module returns the module portion of the address (segments up to and including module names).
func (a *Address) Module() *Address {
	if a == nil || len(a.proto.Segments) == 0 {
		return Empty
	}
	newSegs := make([]*parserpb.Segment, 0, len(a.proto.Segments))
	lastWasModule := false
	for _, seg := range a.proto.Segments {
		switch {
		case lastWasModule:
			newSegs = append(newSegs, seg)
			lastWasModule = false
		case seg.Value == "module":
			newSegs = append(newSegs, seg)
			lastWasModule = true
		default:
			return &Address{proto: &parserpb.Address{Segments: newSegs}}
		}
	}
	return &Address{proto: &parserpb.Address{Segments: newSegs}}
}

// Local returns the non-module portion of the address.
func (a *Address) Local() *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	mod := a.Module()
	return &Address{proto: &parserpb.Address{Segments: a.proto.Segments[len(mod.proto.Segments):]}}
}

// Segments returns the underlying segment slice.
func (a *Address) Segments() []*parserpb.Segment {
	if a == nil || a.proto == nil {
		return nil
	}
	return a.proto.Segments
}

// FromSegments creates an Address from a slice of segments.
func FromSegments(segments []*parserpb.Segment) *Address {
	return &Address{proto: &parserpb.Address{Segments: segments}}
}

// Len returns the number of segments in the address.
func (a *Address) Len() int {
	if a == nil || a.proto == nil {
		return 0
	}
	return len(a.proto.Segments)
}

// IsEmpty returns true if the address has no segments.
func (a *Address) IsEmpty() bool {
	if a == nil || a.proto == nil {
		return true
	}
	return len(a.proto.Segments) == 0
}

// At returns the value of the segment at index i, or empty string if out of bounds.
func (a *Address) At(i int) string {
	if a == nil || a.proto == nil {
		return ""
	}
	if i >= len(a.proto.Segments) {
		return ""
	}
	return a.proto.Segments[i].Value
}

// From returns a new address containing segments from index i onwards.
func (a *Address) From(i int) *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	if i >= len(a.proto.Segments) {
		return &Address{}
	}
	return &Address{proto: &parserpb.Address{Segments: a.proto.Segments[i:]}}
}

// Relative returns the address relative to base, stripping the base prefix if present.
func (a *Address) Relative(base *Address) *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	if base == nil {
		return a
	}
	if !a.StartsWith(base) {
		return a
	}
	rel := &Address{proto: &parserpb.Address{Segments: a.proto.Segments[len(base.proto.Segments):]}}

	if len(base.proto.Segments) > 0 {
		lastBase := base.proto.Segments[len(base.proto.Segments)-1]
		lastSub := a.proto.Segments[len(base.proto.Segments)-1]
		if lastBase.IndexInt == nil && lastBase.IndexString == nil {
			if lastSub.IndexInt != nil || lastSub.IndexString != nil {
				rel.proto.Segments = append([]*parserpb.Segment{{
					IndexInt:    lastSub.IndexInt,
					IndexString: lastSub.IndexString,
				}}, rel.proto.Segments...)
			}
		}
	}

	return rel
}

// New creates a new address from the given parts (strings, ints, or other addresses).
func New(parts ...any) *Address {
	return Empty.Append(parts...)
}

// Parse parses a string representation of an address (e.g., "module.vpc.aws_subnet.public[0]").
func Parse(raw string) (*Address, error) {
	var inQuote bool
	var segments []*parserpb.Segment
	var quoted string
	var bracketed string
	var inBrackets bool
	var current string
	for _, c := range raw {
		if c == '"' {
			inQuote = !inQuote
			continue
		}
		if inQuote {
			quoted += string(c)
			continue
		}
		if c == '[' {
			if inBrackets {
				return nil, fmt.Errorf("unexpected '['")
			}
			if current != "" {
				if number, err := strconv.ParseInt(current, 10, 64); err == nil {
					segments = append(segments, &parserpb.Segment{IndexInt: &number})
				} else {
					segments = append(segments, &parserpb.Segment{Value: current})
				}
			}
			current = ""
			inBrackets = true
			continue
		}
		if c == ']' {
			if !inBrackets {
				return nil, fmt.Errorf("unexpected ']'")
			}

			inBrackets = false
			if quoted != "" {
				if len(segments) > 0 {
					last := segments[len(segments)-1]
					if last.IndexString == nil && last.IndexInt == nil {
						clone := quoted
						last.IndexString = &clone
						quoted = ""
						continue
					}
				}
				segments = append(segments, &parserpb.Segment{
					IndexString: &bracketed,
				})
				quoted = ""
			} else {
				if i, err := strconv.ParseInt(bracketed, 10, 64); err == nil {
					if len(segments) > 0 {
						last := segments[len(segments)-1]
						if last.IndexString == nil && last.IndexInt == nil {
							last.IndexInt = &i
							bracketed = ""
							continue
						}
					}
					segments = append(segments, &parserpb.Segment{
						IndexInt: &i,
					})
				} else {
					return nil, fmt.Errorf("invalid index: %s", bracketed)
				}
				bracketed = ""
			}
			continue
		}
		if inBrackets {
			bracketed += string(c)
			continue
		}
		if c == '.' {
			if current != "" {
				if number, err := strconv.ParseInt(current, 10, 64); err == nil {
					segments = append(segments, &parserpb.Segment{IndexInt: &number})
				} else {
					segments = append(segments, &parserpb.Segment{Value: current})
				}
				current = ""
			}
			continue
		}
		current += string(c)
	}
	if inBrackets {
		return nil, fmt.Errorf("unclosed brackets")
	}
	if inQuote {
		return nil, fmt.Errorf("unclosed quotes")
	}
	if current != "" {
		if number, err := strconv.ParseInt(current, 10, 64); err == nil {
			segments = append(segments, &parserpb.Segment{IndexInt: &number})
		} else {
			segments = append(segments, &parserpb.Segment{Value: current})
		}
	}
	return &Address{proto: &parserpb.Address{Segments: segments}}, nil
}

var stringPool sync.Map

func newInternedString(str string) string {
	if v, ok := stringPool.Load(str); ok {
		return v.(string)
	}
	stringPool.Store(str, str)
	return str
}

// Append returns a new address with the given parts appended.
func (a *Address) Append(parts ...any) *Address {
	if a == nil || a.proto == nil {
		return New(parts...)
	}
	segments := make([]*parserpb.Segment, len(a.proto.Segments), len(a.proto.Segments)+len(parts))
	if len(a.proto.Segments) > 0 {
		copy(segments, a.proto.Segments)
	}
	for i := range parts {
		switch v := parts[i].(type) {
		case *parserpb.Address:
			if v != nil {
				segments = append(segments, v.Segments...)
			}
		case *Address:
			if v != nil {
				segments = append(segments, v.proto.Segments...)
			}
		case []string:
			for _, part := range v {
				segments = append(segments, &parserpb.Segment{Value: newInternedString(part)})
			}
		case string:
			segments = append(segments, &parserpb.Segment{Value: newInternedString(v)})
		case int:
			i64 := int64(v)
			if len(a.proto.Segments) == 0 || a.proto.Segments[len(a.proto.Segments)-1].IndexInt != nil {
				segments = append(segments, &parserpb.Segment{
					IndexInt: &i64,
				})
			} else {
				segments[len(segments)-1].IndexInt = &i64
			}
		case int64:
			if len(a.proto.Segments) == 0 || a.proto.Segments[len(a.proto.Segments)-1].IndexInt != nil {
				segments = append(segments, &parserpb.Segment{
					IndexInt: &v,
				})
			} else {
				segments[len(segments)-1].IndexInt = &v
			}
		case float64:
			conv := int64(v)
			if len(a.proto.Segments) == 0 || a.proto.Segments[len(a.proto.Segments)-1].IndexInt != nil {
				segments = append(segments, &parserpb.Segment{
					IndexInt: &conv,
				})
			} else {
				segments[len(segments)-1].IndexInt = &conv
			}
		default:
			segments = append(segments, &parserpb.Segment{Value: newInternedString(fmt.Sprintf("%v", v))})
		}
	}
	return &Address{proto: &parserpb.Address{Segments: segments}}
}

// Equal returns true if the two addresses are equal.
func (a *Address) Equal(b *Address) bool {
	if (a == nil || a.Len() == 0) && (b == nil || b.Len() == 0) {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a.proto.Segments) != len(b.proto.Segments) {
		return false
	}
	iterA := NewIterator(a)
	iterB := NewIterator(b)
	for {
		if iterA.IsDone() && iterB.IsDone() {
			return true
		}
		if iterA.IsDone() != iterB.IsDone() {
			return false
		}
		a1, a2 := iterA.PopKey()
		b1, b2 := iterB.PopKey()
		if (a1 == nil) != (b1 == nil) {
			return false
		}
		if (a2 == nil) != (b2 == nil) {
			return false
		}
		if a1 != nil {
			if *a1 != *b1 {
				return false
			}
		}
		if a2 != nil {
			if *a2 != *b2 {
				return false
			}
		}
	}
}

// Clone returns a deep copy of the address.
func (a *Address) Clone() *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	clone := make([]*parserpb.Segment, len(a.proto.Segments))
	for i, seg := range a.proto.Segments {
		clone[i] = &parserpb.Segment{
			Value:       seg.Value,
			IndexString: seg.IndexString,
			IndexInt:    seg.IndexInt,
		}
	}
	return &Address{proto: &parserpb.Address{Segments: clone}}
}

// ToGraph returns a copy of the address with all indices removed.
func (a *Address) ToGraph() *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	var segments []*parserpb.Segment
	for _, seg := range a.proto.Segments {
		if seg.Value != "" {
			segments = append(segments, &parserpb.Segment{Value: seg.Value})
		}
	}
	return &Address{proto: &parserpb.Address{Segments: segments}}
}

// Truncate returns a new address containing only the first n segments.
func (a *Address) Truncate(length int) *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	if length > len(a.proto.Segments) {
		length = len(a.proto.Segments)
	}
	return &Address{
		proto: &parserpb.Address{Segments: a.proto.Segments[:length]},
	}
}

// CreateChild returns a new address with the given names appended as segments.
func (a *Address) CreateChild(name ...string) *Address {
	if a == nil || a.proto == nil {
		parts := make([]any, 0, len(name))
		for _, n := range name {
			parts = append(parts, n)
		}
		return New(parts...)

	}
	segments := make([]*parserpb.Segment, len(a.proto.Segments)+len(name))
	copy(segments, a.proto.Segments)
	for i, n := range name {
		segments[len(a.proto.Segments)+i] = &parserpb.Segment{Value: n}
	}
	return &Address{proto: &parserpb.Address{Segments: segments}}
}

// CreateStringIndexedChild returns a new address with a string index appended to the last segment.
func (a *Address) CreateStringIndexedChild(index string) *Address {
	if a.Len() == 0 {
		return &Address{proto: &parserpb.Address{Segments: []*parserpb.Segment{{IndexString: &index}}}}
	}
	clone := a.Clone()
	last := clone.proto.Segments[len(clone.proto.Segments)-1]
	if last.IndexString != nil || last.IndexInt != nil {
		clone.proto.Segments = append(clone.proto.Segments, &parserpb.Segment{})
	}
	index = newInternedString(index)
	clone.proto.Segments[len(clone.proto.Segments)-1].IndexString = &index
	return clone
}

// CreateIntIndexedChild returns a new address with an integer index appended to the last segment.
func (a *Address) CreateIntIndexedChild(index int64) *Address {
	if a.Len() == 0 {
		return &Address{proto: &parserpb.Address{Segments: []*parserpb.Segment{{IndexInt: &index}}}}
	}
	clone := a.Clone()
	last := clone.proto.Segments[len(clone.proto.Segments)-1]
	if last.IndexString != nil || last.IndexInt != nil {
		clone.proto.Segments = append(clone.proto.Segments, &parserpb.Segment{})
	}
	clone.proto.Segments[len(clone.proto.Segments)-1].IndexInt = &index
	return clone
}

// LastIntIndex returns the integer index of the last segment, or -1 if none.
func (a *Address) LastIntIndex() int64 {
	if a == nil || len(a.proto.Segments) == 0 {
		return -1
	}
	last := a.proto.Segments[len(a.proto.Segments)-1]
	if last.IndexInt != nil {
		return *last.IndexInt
	}
	return -1
}

// Last returns the value of the last segment, or empty string if empty.
func (a *Address) Last() string {
	if a == nil || a.proto == nil {
		return ""
	}
	if len(a.proto.Segments) == 0 {
		return ""
	}
	return a.proto.Segments[len(a.proto.Segments)-1].Value
}

// StartsWith returns true if the address starts with the given prefix.
func (a *Address) StartsWith(other *Address) bool {
	return a.Truncate(other.Len()).Equal(other) || a.Truncate(other.Len()).StripIndex().Equal(other)
}

// StripIndex returns a copy with the index removed from the last segment.
func (a *Address) StripIndex() *Address {
	if a == nil || a.proto == nil {
		return Empty
	}
	clone := a.Clone()
	if len(clone.proto.Segments) >= 1 {
		clone.proto.Segments[len(clone.proto.Segments)-1].IndexInt = nil
		clone.proto.Segments[len(clone.proto.Segments)-1].IndexString = nil
	}
	return clone
}

// Hash returns a hex-encoded FNV-64 hash of the address string.
func (a *Address) Hash() string {
	h := fnv.New64()
	_, _ = h.Write([]byte(a.String()))
	return hex.EncodeToString(h.Sum(nil))
}

// String returns the string representation of the address (e.g., "module.vpc.aws_subnet.public[0]").
func (a *Address) String() string {
	if a == nil || a.proto == nil {
		return ""
	}
	if len(a.proto.Segments) == 0 {
		return ""
	}
	a.strMu.RLock()
	if a.str != nil {
		defer a.strMu.RUnlock()
		return *a.str
	}
	a.strMu.RUnlock()
	a.strMu.Lock()
	defer a.strMu.Unlock()

	// double check nothing changed while we were waiting
	if a.str != nil {
		return *a.str
	}

	strs := make([]string, 0, len(a.proto.Segments))
	for _, seg := range a.proto.Segments {
		switch {
		case (seg.IndexInt == nil && seg.IndexString == nil):
			strs = append(strs, seg.Value)
		case seg.IndexInt != nil:
			if seg.Value == "" && len(strs) > 0 {
				strs[len(strs)-1] += "[" + strconv.Itoa(int(*seg.IndexInt)) + "]"
			} else {
				strs = append(strs, seg.Value+"["+strconv.Itoa(int(*seg.IndexInt))+"]")
			}
		default:
			if seg.Value == "" && len(strs) > 0 {
				strs[len(strs)-1] += `["` + *seg.IndexString + `"]`
			} else {
				strs = append(strs, seg.Value+`["`+*seg.IndexString+`"]`)
			}
		}
	}
	combined := strings.Join(strs, ".")
	a.str = &combined

	return *a.str
}

// RelativeLooksLikeIntIndex returns true if the first segment looks like an integer index.
func (a *Address) RelativeLooksLikeIntIndex() bool {
	if a == nil || len(a.proto.Segments) == 0 {
		return false
	}
	first := a.proto.Segments[0]
	if first.IndexInt != nil {
		return true
	}
	_, err := strconv.Atoi(first.Value)
	return err == nil
}
