package resource

import (
	"encoding/hex"
	"hash/fnv"
	"reflect"
	"sort"
	"strings"

	"github.com/infracost/go-proto/pkg/address"
	"github.com/infracost/go-proto/pkg/flag"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/infracost/proto/gen/go/infracost/provider"
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
	IsProduction           bool
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
	Key       value.String `tree:"key"`
	Value     value.String `tree:"value"`
	IsDefault bool         `tree:"-"` // not needed for our purposes
}

func (t Tags) Get(name string) (value.String, bool) {
	for _, tag := range t {
		if tag.Key.Equal(name) {
			return tag.Value, true
		}
	}
	return value.EmptyString, false
}

func (t Tags) GetDefaults() Tags {
	var output Tags
	for _, tag := range t {
		if tag.IsDefault {
			output = append(output, tag)
		}
	}
	return output
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

func (t *Tags) Set(k, v value.String, isDefault bool) {
	for i, existing := range *t {
		if existing.Key.String() == k.String() {
			if existing.IsDefault || !isDefault {
				(*t)[i].Value = v
				(*t)[i].IsDefault = isDefault
			}
			return
		}
	}
	*t = append(*t, Tag{
		Key:       k,
		Value:     v,
		IsDefault: isDefault,
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

// ToProviderResource takes the wrapping Implementation (e.g. *ec2.Instance)
// rather than being a method on *Resource so the deep-checksum walker can
// reflect over the full struct — including wrapper-struct fields like
// Relationships — instead of just the embedded base.
func ToProviderResource(impl Implementation) *provider.Resource {
	base := impl.GetBase()
	c := BaseToProviderResource(base)
	c.Metadata.DeepChecksum = CalculateDeepChecksum(impl)
	return c
}

// ProtoToProviderResource converts a *prototree.Resource directly to a
// *provider.Resource without first decoding into the intermediate Go
// Resource struct. Use this when you only have the proto representation
// in hand and don't need the typed helpers attached to the Go Resource.
//
// DeepChecksum on the returned message is set to FullChecksum — the
// deep-checksum walker requires a typed Implementation, so callers that
// need a true deep checksum should go through ToProviderResource.
func ProtoToProviderResource(res *prototree.Resource) *provider.Resource {
	if res == nil {
		return nil
	}

	resourceMetadata := provider.ResourceMetadata{
		DefaultTagsChecksum:    defaultTagsChecksumFromProto(res.Tags),
		BasicChecksum:          res.BasicChecksum,
		AttributeValueChecksum: res.FullChecksum,
		DeepChecksum:           res.FullChecksum, // overridden by consumer
	}

	var (
		callStack    *parser.CallStack
		resourceType string
		name         string
		providerLink string
	)
	if def := res.Definition; def != nil {
		callStack = def.CallStack
		resourceType = def.ResourceType
		name = address.FromProto(def.Address).String()

		if def.CallStack != nil {
			var fullAddressSize int
			resourceMetadata.ModuleCalls = make([]*provider.ModuleCall, 0, len(def.CallStack.Frames))
			for _, frame := range def.CallStack.Frames {
				resourceMetadata.ModuleCalls = append(resourceMetadata.ModuleCalls, &provider.ModuleCall{
					Filename:       frame.SourceRange.Filename,
					StartLine:      frame.SourceRange.StartLine,
					EndLine:        frame.SourceRange.EndLine,
					DefinitionName: address.FromProto(frame.Address).From(fullAddressSize).String(),
				})
				fullAddressSize = address.FromProto(frame.Address).Len()
			}
		}

		if srcRange := def.Source; srcRange != nil {
			resourceMetadata.StartLine = srcRange.StartLine
			resourceMetadata.EndLine = srcRange.EndLine
			resourceMetadata.Filename = srcRange.Filename
			resourceMetadata.ModuleCalls = nil
		}

		if pc := def.ProviderConfiguration; pc != nil && pc.Source != nil {
			providerLink = pc.Source.Filename
		}
	}

	outputResource := &provider.Resource{
		Id:                  res.Id,
		Metadata:            &resourceMetadata,
		Type:                resourceType,
		ProviderLink:        providerLink,
		Name:                name,
		IsSupported:         false, // overridden by consumer
		IsFree:              res.IsFree,
		IsProviderSupported: true,                      // overridden by consumer
		Costs:               &provider.ResourceCosts{}, // overridden by consumer
		ChildResources:      nil,
		Tagging: &provider.Tagging{
			Tags:                make([]*provider.Tag, 0, len(res.Tags)),
			SupportsDefaultTags: res.SupportsDefaultTags,
			SupportsTags:        res.SupportsTags,
			PropagationProblems: convertTagPropagationProblems(res.TagPropagationProblems),
		},
		Region:    res.Region,
		Action:    provider.ResourceAction_RESOURCE_ACTION_UNSPECIFIED,
		CallStack: callStack,
	}

	// Match BaseToProviderResource's observable tag ordering: tags come out
	// sorted by key (a side effect of DefaultChecksum sorting Tags in place).
	sortedTags := make([]*prototree.Tag, len(res.Tags))
	copy(sortedTags, res.Tags)
	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].GetKey().GetStringValue() < sortedTags[j].GetKey().GetStringValue()
	})
	for _, tag := range sortedTags {
		outputResource.Tagging.Tags = append(outputResource.Tagging.Tags, &provider.Tag{
			Key:              tag.GetKey().GetStringValue(),
			Value:            tag.GetValue().GetStringValue(),
			IsDefault:        tag.IsDefault,
			IsKeySynthetic:   flag.Flags(tag.GetKey().GetFlags()).IsSynthetic(),
			IsValueSynthetic: flag.Flags(tag.GetValue().GetFlags()).IsSynthetic(),
		})
	}

	return outputResource
}

func defaultTagsChecksumFromProto(tags []*prototree.Tag) string {
	sorted := make([]*prototree.Tag, len(tags))
	copy(sorted, tags)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetKey().GetStringValue() < sorted[j].GetKey().GetStringValue()
	})

	var b strings.Builder
	for _, tag := range sorted {
		if !tag.IsDefault {
			continue
		}
		b.WriteString(tag.GetKey().GetStringValue())
		b.WriteString("=")
		b.WriteString(tag.GetValue().GetStringValue())
		b.WriteString("\n")
	}
	return hash(b.String())
}

// BaseToProviderResource converts the base Resource struct to a provider.Resource proto message.
func BaseToProviderResource(base *Resource) *provider.Resource {
	resourceMetadata := provider.ResourceMetadata{
		DefaultTagsChecksum:    base.Tags.DefaultChecksum(),
		BasicChecksum:          base.BasicChecksum,
		AttributeValueChecksum: base.FullChecksum,
		DeepChecksum:           base.FullChecksum, // overridden by consumer or by impl hashing
	}

	var fullAddressSize int
	if base.Definition.CallStack != nil {
		resourceMetadata.ModuleCalls = make([]*provider.ModuleCall, 0, len(base.Definition.CallStack.Frames))
		for _, frame := range base.Definition.CallStack.Frames {
			resourceMetadata.ModuleCalls = append(resourceMetadata.ModuleCalls, &provider.ModuleCall{
				Filename:       frame.SourceRange.Filename,
				StartLine:      frame.SourceRange.StartLine,
				EndLine:        frame.SourceRange.EndLine,
				DefinitionName: address.FromProto(frame.Address).From(fullAddressSize).String(),
			})
			fullAddressSize = address.FromProto(frame.Address).Len()
		}
	}

	if srcRange := base.Definition.SourceRange; srcRange != nil {
		resourceMetadata.StartLine = srcRange.StartLine
		resourceMetadata.EndLine = srcRange.EndLine
		resourceMetadata.Filename = srcRange.Filename
		resourceMetadata.ModuleCalls = nil
	}

	providerLink := ""
	if pc := base.Definition.ProviderConfiguration; pc != nil && pc.Source != nil {
		providerLink = pc.Source.Filename
	}

	outputResource := &provider.Resource{
		Id:                  base.ID,
		Metadata:            &resourceMetadata,
		Type:                base.Definition.ResourceType,
		ProviderLink:        providerLink,
		Name:                base.Definition.Address.String(),
		IsSupported:         false, // overridden by consumer
		IsFree:              base.IsFree,
		IsProviderSupported: true,                      // overridden by consumer
		Costs:               &provider.ResourceCosts{}, // overridden by consumer
		ChildResources:      nil,
		Tagging: &provider.Tagging{
			Tags:                make([]*provider.Tag, 0, len(base.Tags)), // populated later on...
			SupportsDefaultTags: base.SupportsDefaultTags,
			SupportsTags:        base.SupportsTags,
			PropagationProblems: convertTagPropagationProblems(base.TagPropagationProblems),
		},
		Region:    base.Region,
		Action:    provider.ResourceAction_RESOURCE_ACTION_UNSPECIFIED,
		CallStack: base.Definition.CallStack,
	}

	for _, tag := range base.Tags {
		outputResource.Tagging.Tags = append(outputResource.Tagging.Tags, &provider.Tag{
			Key:              tag.Key.Value(),
			Value:            tag.Value.Value(),
			IsDefault:        tag.IsDefault,
			IsKeySynthetic:   tag.Key.Flags().IsSynthetic(),
			IsValueSynthetic: tag.Value.Flags().IsSynthetic(),
		})
	}

	return outputResource
}

func convertTagPropagationProblems(problems []*prototree.TagPropagationProblem) []*provider.TagPropagationProblem {
	protoProblems := make([]*provider.TagPropagationProblem, 0, len(problems))
	for _, problem := range problems {
		protoProblems = append(protoProblems, &provider.TagPropagationProblem{
			ActualValue:  problem.ActualValue,
			ValidValues:  problem.ValidValues,
			Attribute:    problem.Attribute,
			TagRecipient: problem.TagRecipient,
			AffectedTags: problem.AffectedTags,
		})
	}
	return protoProblems
}

// maxChecksumDepth limits how deep the recursive checksum calculation will traverse
// linked resources. This prevents excessive computation on deeply nested resource
// graphs and guards against potential circular references. A depth of 5 is sufficient
// to capture meaningful relationships while keeping the computation bounded.
const maxChecksumDepth = 5

func CalculateDeepChecksum(impl Implementation) string {
	return calculateDeepChecksum(impl, 0)
}

type Implementation interface {
	GetBase() *Resource
}

func calculateDeepChecksum(impl Implementation, depth int) (checksum string) {
	if impl == nil {
		return ""
	}

	// in the event of a panic, use the base checksum, it's always going to be a nil resource,
	// or an incorrectly implemented resource with no checksum
	defer func() {
		if p := recover(); p != nil {
			checksum = impl.GetBase().FullChecksum
		}
	}()

	checksum = impl.GetBase().FullChecksum

	depth++
	if depth >= maxChecksumDepth {
		return
	}

	val := reflect.ValueOf(impl)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return checksum
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return checksum
	}

	checksums := []string{checksum}
	checksums = append(checksums, collectImplementationChecksums(val, depth)...)
	return hash(strings.Join(checksums, "-"))
}

// collectImplementationChecksums walks an arbitrary value looking for nested
// Implementations and returns their deep checksums. Traversing non-Implementation
// containers (plain structs, pointers-to-plain-structs, slices, maps) does not
// consume depth — only calculateDeepChecksum does, when it actually recurses
// into an Implementation. This is what lets Instance.Relationships.Whatever be
// picked up even though Relationships itself isn't an Implementation.
func collectImplementationChecksums(val reflect.Value, depth int) []string {
	implType := reflect.TypeFor[Implementation]()

	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		if val.Type().Implements(implType) {
			return []string{calculateDeepChecksum(val.Interface().(Implementation), depth)}
		}
		return collectImplementationChecksums(val.Elem(), depth)

	case reflect.Struct:
		var out []string
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanInterface() {
				continue
			}
			out = append(out, collectImplementationChecksums(field, depth)...)
		}
		return out

	case reflect.Slice, reflect.Array:
		var elems []string
		for i := 0; i < val.Len(); i++ {
			elems = append(elems, collectImplementationChecksums(val.Index(i), depth)...)
		}
		sort.Strings(elems)
		return elems

	case reflect.Map:
		var elems []string
		iter := val.MapRange()
		for iter.Next() {
			elems = append(elems, collectImplementationChecksums(iter.Value(), depth)...)
		}
		sort.Strings(elems)
		return elems
	}
	return nil
}
