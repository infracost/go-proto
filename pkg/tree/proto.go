package tree

import (
	"fmt"
	"reflect"

	"github.com/infracost/go-proto/pkg/address"
	"github.com/infracost/go-proto/pkg/tree/resource"
	prototree "github.com/infracost/proto/gen/go/infracost/tree"
)

var resourceType = reflect.TypeFor[resource.Resource]()

func (t *Tree) ToProto() (*prototree.Tree, error) {
	pt := &prototree.Tree{
		Providers: make(map[string]*prototree.Provider),
	}

	tv := reflect.ValueOf(t).Elem()
	tt := tv.Type()

	for i := range tt.NumField() {
		providerTag := tt.Field(i).Tag.Get("tree")
		if providerTag == "" || providerTag == "-" {
			continue
		}

		providerVal := tv.Field(i)

		protoProvider := &prototree.Provider{
			Services: make(map[string]*prototree.Service),
		}
		providerType := providerVal.Type()

		for j := range providerType.NumField() {
			serviceTag := providerType.Field(j).Tag.Get("tree")
			if serviceTag == "" || serviceTag == "-" {
				continue
			}

			service := &prototree.Service{}
			serviceVal := providerVal.Field(j)
			serviceType := serviceVal.Type()

			for k := range serviceType.NumField() {
				resourceTag := serviceType.Field(k).Tag.Get("tree")
				if resourceTag == "" || resourceTag == "-" {
					continue
				}

				sliceVal := serviceVal.Field(k)
				for l := range sliceVal.Len() {
					item := sliceVal.Index(l)
					attrs := StructToValueObject(item.Interface())
					base := baseResource(item)

					var providerConfig *prototree.ProviderConfiguration
					if config := base.Definition.ProviderConfiguration; config != nil {
						providerConfig = &prototree.ProviderConfiguration{
							Source:             config.Source,
							VersionConstraints: config.VersionConstraints,
						}
					}

					res := convertResourceToProto(
						&base,
						resourceTag,
						attrs,
						providerConfig,
					)
					service.Resources = append(service.Resources, res)
				}
			}

			protoProvider.Services[serviceTag] = service
		}

		pt.Providers[providerTag] = protoProvider
	}

	pt.UnsupportedResources = make([]*prototree.Resource, len(t.UnsupportedResources))
	for i, u := range t.UnsupportedResources {
		pt.UnsupportedResources[i] = convertResourceToProto(
			u,
			"",
			nil,
			nil,
		)
	}

	return pt, nil
}

func FromProto(p *prototree.Tree) (*Tree, error) {
	t := &Tree{}

	tv := reflect.ValueOf(t).Elem()
	tt := tv.Type()

	for i := range tt.NumField() {
		providerTag := tt.Field(i).Tag.Get("tree")
		if providerTag == "" {
			continue
		}

		provider, ok := p.Providers[providerTag]
		if !ok {
			continue
		}

		providerVal := tv.Field(i)
		providerType := providerVal.Type()

		for j := range providerType.NumField() {
			serviceTag := providerType.Field(j).Tag.Get("tree")
			if serviceTag == "" {
				continue
			}

			service, ok := provider.Services[serviceTag]
			if !ok {
				continue
			}

			serviceVal := providerVal.Field(j)
			serviceType := serviceVal.Type()

			for k := range serviceType.NumField() {
				resourceTag := serviceType.Field(k).Tag.Get("tree")
				if resourceTag == "" {
					continue
				}

				sliceField := serviceVal.Field(k)
				elemType := sliceField.Type().Elem()

				for _, res := range service.Resources {
					if res.Type != resourceTag {
						continue
					}

					item := reflect.New(elemType).Elem()
					setBaseResource(item, res)
					if res.Attributes != nil {
						ValueObjectToStruct(res.Attributes, item.Addr().Interface())
					}
					sliceField.Set(reflect.Append(sliceField, item))
				}
			}
		}
	}

	if len(p.Providers) > 0 {
		for name := range p.Providers {
			found := false
			for i := range tt.NumField() {
				if tt.Field(i).Tag.Get("tree") == name {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("unmapped provider: %s", name)
			}
		}
	}

	t.UnsupportedResources = make([]*resource.Resource, len(p.UnsupportedResources))
	for i, p := range p.UnsupportedResources {
		t.UnsupportedResources[i] = convertResourceFromProto(p, nil)
	}

	return t, nil
}

func baseResource(v reflect.Value) resource.Resource {
	f := v.FieldByName("Resource")
	if !f.IsValid() || f.Type() != resourceType {
		return resource.Resource{}
	}
	return f.Interface().(resource.Resource)
}

func setBaseResource(v reflect.Value, res *prototree.Resource) {
	f := v.FieldByName("Resource")
	if !f.IsValid() || f.Type() != resourceType {
		return
	}

	var providerConfig *resource.ProviderConfiguration
	if res.Definition != nil {
		if c := res.Definition.ProviderConfiguration; c != nil {
			providerConfig = &resource.ProviderConfiguration{
				Source:             c.Source,
				VersionConstraints: c.VersionConstraints,
			}
		}
	}

	f.Set(reflect.ValueOf(convertResourceFromProto(res, providerConfig)).Elem())
}

func convertResourceFromProto(res *prototree.Resource, providerConfig *resource.ProviderConfiguration) *resource.Resource {
	r := &resource.Resource{
		ID:                     res.Id,
		Region:                 res.Region,
		IsDebug:                res.IsDebug,
		IsChild:                res.IsChild,
		IsFree:                 res.IsFree,
		RegionIsSynthetic:      res.RegionIsSynthetic,
		SupportsTags:           res.SupportsTags,
		SupportsDefaultTags:    res.SupportsDefaultTags,
		BasicChecksum:          res.BasicChecksum,
		FullChecksum:           res.FullChecksum,
		Flags:                  res.Flags,
		Tags:                   resource.TagsFromProto(res.Tags),
		TagPropagationProblems: res.TagPropagationProblems,
	}
	if res.Definition != nil {
		r.Definition = resource.Definition{
			CallStack:             res.Definition.CallStack,
			ProviderConfiguration: providerConfig,
			SourceRange:           res.Definition.Source,
			ResourceType:          res.Definition.ResourceType,
			Address:               address.FromProto(res.Definition.Address),
			RawStringAttributes:   res.Definition.RawStringAttributes,
		}
	}
	return r
}

func convertResourceToProto(base *resource.Resource, resourceTag string, attrs *prototree.ValueObject, providerConfig *prototree.ProviderConfiguration) *prototree.Resource {
	return &prototree.Resource{
		Attributes:             attrs,
		Id:                     base.ID,
		Type:                   resourceTag,
		Region:                 base.Region,
		IsDebug:                base.IsDebug,
		IsChild:                base.IsChild,
		IsFree:                 base.IsFree,
		RegionIsSynthetic:      base.RegionIsSynthetic,
		SupportsTags:           base.SupportsTags,
		SupportsDefaultTags:    base.SupportsDefaultTags,
		BasicChecksum:          base.BasicChecksum,
		FullChecksum:           base.FullChecksum,
		Flags:                  base.Flags,
		Tags:                   base.Tags.ToProto(),
		TagPropagationProblems: base.TagPropagationProblems,
		Definition: &prototree.Definition{
			CallStack:             base.Definition.CallStack,
			Source:                base.Definition.SourceRange,
			ResourceType:          base.Definition.ResourceType,
			Address:               base.Definition.Address.ToProto(),
			RawStringAttributes:   base.Definition.RawStringAttributes,
			ProviderConfiguration: providerConfig,
		},
	}
}
