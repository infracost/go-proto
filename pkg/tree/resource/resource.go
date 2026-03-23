package resource

import (
	"encoding/hex"
	"hash/fnv"
	"sort"
	"strings"

	"github.com/infracost/go-proto/pkg/address"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/infracost/proto/gen/go/infracost/parser"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
)

// Resource is a base struct embedded by all resource types (e.g. ec2.Instance).
// It contains the common properties shared across all resources in the tree.
type Resource struct {
	ID                     string
	Region                 string
	RegionIsSynthetic      bool
	IsDebug                bool
	IsChild                bool
	IsFree                 bool
	SupportsTags           bool
	SupportsDefaultTags    bool
	BasicChecksum          string
	FullChecksum           string
	Flags                  uint64
	Tags                   Tags
	Definition             Definition
	TagPropagationProblems []*prototree.TagPropagationProblem
}

func (r *Resource) GetBase() *Resource {
	return r
}

type ProviderConfiguration struct {
	Source             *parser.SourceRange
	VersionConstraints string
}

type Definition struct {
	CallStack             *parser.CallStack
	ProviderConfiguration *ProviderConfiguration
	SourceRange           *parser.SourceRange
	ResourceType          string
	Address               *address.Address
	RawStringAttributes   map[string]string
}

type Tags []Tag

type Tag struct {
	Key       value.String
	Value     value.String
	IsDefault bool
}

func (t Tags) Get(name string) (value.String, bool) {
	for _, tag := range t {
		if tag.Key.Equals(name) {
			return tag.Value, true
		}
	}
	return value.EmptyString, false
}

func (t Tags) DefaultChecksum() string {
	sort.Slice(t, func(i, j int) bool {
		is := t[i].Key.Value()
		js := t[j].Key.Value()
		return is < js
	})
	var b strings.Builder
	for _, tag := range t {
		if !tag.IsDefault {
			continue
		}
		ks := tag.Key.Value()
		vs := tag.Value.Value()
		b.WriteString(ks)
		b.WriteString("=")
		b.WriteString(vs)
		b.WriteString("\n")
	}
	return hash(b.String())
}

func (t *Tags) Set(k, v value.String, overwrite bool) {
	for i, existing := range *t {
		if existing.Key.String() == k.String() {
			if overwrite {
				(*t)[i].Value = v
			}
			return
		}
	}
	*t = append(*t, Tag{
		Key:       k,
		Value:     v,
		IsDefault: false,
	})
}

func (t Tags) ToProto() []*prototree.Tag {
	output := make([]*prototree.Tag, 0, len(t))
	for _, tag := range t {
		output = append(output, &prototree.Tag{
			Key:       tag.Key.ToProto(),
			Value:     tag.Value.ToProto(),
			IsDefault: tag.IsDefault,
		})
	}
	return output
}

func TagsFromProto(input []*prototree.Tag) Tags {
	output := make(Tags, 0, len(input))
	for _, tag := range input {
		output = append(output, Tag{
			Key:       value.FromProto[string](tag.Key),
			Value:     value.FromProto[string](tag.Value),
			IsDefault: tag.IsDefault,
		})
	}
	return output
}

func hash(s string) string {
	h := fnv.New64()
	_, _ = h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
