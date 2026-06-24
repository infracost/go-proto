package tree

import (
	"fmt"
	"reflect"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/azure"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/kubernetes"
	"github.com/infracost/go-proto/pkg/tree/resource"
)

var resourceIfaceType = reflect.TypeFor[Resource]()

type Resource interface {
	GetBase() *resource.Resource
}

type Tree struct {
	AWS                  aws.AWS               `tree:"aws"`
	Azure                azure.Azure           `tree:"azure"`
	Google               google.Google         `tree:"google"`
	Kubernetes           kubernetes.Kubernetes `tree:"kubernetes"`
	UnsupportedResources []*resource.Resource  `tree:"-"` // these get handled as a special case
}

// ToResources returns every struct in the tree that embeds resource.Resource.
func (t *Tree) ToResources(includeUnsupported bool) []Resource {
	var output []Resource
	collectResources(reflect.ValueOf(t).Elem(), &output)
	if includeUnsupported {
		for _, res := range t.UnsupportedResources {
			output = append(output, res)
		}
	}
	return output
}

func collectResources(v reflect.Value, out *[]Resource) {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			collectResources(v.Elem(), out)
		}
	case reflect.Struct:
		if _, ok := v.Type().FieldByName("Resource"); ok {
			f := v.FieldByName("Resource")
			if f.IsValid() && f.Type() == resourceType && reflect.PointerTo(v.Type()).Implements(resourceIfaceType) {
				*out = append(*out, v.Addr().Interface().(Resource))
				return
			}
		}
		t := v.Type()
		for i := range t.NumField() {
			tag := t.Field(i).Tag.Get("tree")
			if tag == "" || tag == "-" {
				continue
			}
			collectResources(v.Field(i), out)
		}
	case reflect.Slice:
		for i := range v.Len() {
			collectResources(v.Index(i), out)
		}
	}
}

// ModifyResource finds the equivalent resource of _target_ in a cloned tree and modifies it with the supplied function
// Note that _target_ MUST be a pointer for this to compile
func ModifyResource[T interface {
	Resource
	~*E
}, E any](t *Tree, target T, modify func(t T)) error {
	for _, resource := range t.ToResources(true) {
		if match, ok := resource.(T); ok && match.GetBase().Definition.Address.Equal(target.GetBase().Definition.Address) {
			modify(match)
			return nil
		}
	}
	return fmt.Errorf("failed to find resource %s", target.GetBase().Definition.Address.String())
}

type postProcessor interface {
	PostProcess()
}

func (t *Tree) PostProcess() {
	postProcess(reflect.ValueOf(t).Elem())
}

func postProcess(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			postProcess(v.Elem())
		}
	case reflect.Struct:
		t := v.Type()
		for i := range t.NumField() {
			tag := t.Field(i).Tag.Get("tree")
			if tag == "" || tag == "-" {
				continue
			}
			f := v.Field(i)
			if f.Kind() == reflect.Struct && f.CanAddr() {
				if pp, ok := f.Addr().Interface().(postProcessor); ok {
					pp.PostProcess()
				}
			}
			postProcess(f)
		}
	case reflect.Slice:
		for i := range v.Len() {
			postProcess(v.Index(i))
		}
	}
}
