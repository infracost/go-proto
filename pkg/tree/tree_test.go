package tree

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToResources(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
					{Resource: resource.Resource{ID: "i-222"}},
				},
			},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	inst0, ok := resources[0].(*ec2.Instance)
	require.True(t, ok)
	assert.Equal(t, "i-111", inst0.ID)

	inst1, ok := resources[1].(*ec2.Instance)
	require.True(t, ok)
	assert.Equal(t, "i-222", inst1.ID)
}

func TestToResourcesEmpty(t *testing.T) {
	tree := &Tree{}
	resources := tree.ToResources()
	assert.Empty(t, resources)
}

func TestToResourcesUnsupported(t *testing.T) {
	tree := &Tree{
		UnsupportedResources: []*resource.Resource{
			{ID: "unsupported-1", Region: "us-west-2"},
			{ID: "unsupported-2", Region: "eu-west-1"},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	assert.Equal(t, "unsupported-1", resources[0].GetBase().ID)
	assert.Equal(t, "us-west-2", resources[0].GetBase().Region)
	assert.Equal(t, "unsupported-2", resources[1].GetBase().ID)
	assert.Equal(t, "eu-west-1", resources[1].GetBase().Region)
}

func TestPostProcess_CallsChildPostProcess(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
				},
				InstanceStates: []ec2.InstanceStateMapping{
					{InstanceID: value.New("i-111", 0, "", nil)},
				},
			},
		},
	}

	tree.PostProcess()

	// EC2.PostProcess should have linked the instance state
	require.NotNil(t, tree.AWS.EC2.Instances[0].Relationships.InstanceState)
	assert.Equal(t, &tree.AWS.EC2.InstanceStates[0], tree.AWS.EC2.Instances[0].Relationships.InstanceState)
}

func TestPostProcess_EmptyTree(t *testing.T) {
	tree := &Tree{}
	// should not panic
	tree.PostProcess()
}

func TestToResourcesMixed(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-111"}},
				},
			},
		},
		UnsupportedResources: []*resource.Resource{
			{ID: "unsupported-1"},
		},
	}

	resources := tree.ToResources()
	require.Len(t, resources, 2)

	assert.Equal(t, "i-111", resources[0].GetBase().ID)
	assert.Equal(t, "unsupported-1", resources[1].GetBase().ID)
}

func TestAllFieldsHaveTreeTag(t *testing.T) {
	var problems []string
	seen := make(map[reflect.Type]bool)
	checkTreeTags(reflect.TypeOf(Tree{}), "Tree", &problems, seen)
	for _, p := range problems {
		t.Error(p)
	}
}

const treePkgPrefix = "github.com/infracost/go-proto/pkg/tree"

var valuePkgPath = reflect.TypeOf(value.Value[string]{}).PkgPath()

func isTreePackage(pkgPath string) bool {
	return pkgPath == treePkgPrefix || len(pkgPath) > len(treePkgPrefix) && pkgPath[len(treePkgPrefix)] == '/'
}

func deref(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func checkTreeTags(typ reflect.Type, path string, problems *[]string, seen map[reflect.Type]bool) {
	typ = deref(typ)
	if typ.Kind() != reflect.Struct || seen[typ] {
		return
	}
	seen[typ] = true

	tagValues := make(map[string]string) // tag value -> field name
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldPath := fmt.Sprintf("%s.%s", path, field.Name)

		tag, hasTag := field.Tag.Lookup("tree")
		if !hasTag {
			*problems = append(*problems, fmt.Sprintf("missing tree tag: %s", fieldPath))
			continue
		}
		if tag == "-" {
			continue
		}

		// check for duplicate tag values within the same struct
		if prev, exists := tagValues[tag]; exists {
			*problems = append(*problems, fmt.Sprintf("duplicate tree tag %q: %s.%s and %s.%s", tag, path, prev, path, field.Name))
		} else {
			tagValues[tag] = field.Name
		}

		// resolve element type for slices/pointers
		ft := deref(field.Type)
		if ft.Kind() == reflect.Slice {
			ft = deref(ft.Elem())
		}

		// stop at value package types
		if ft.PkgPath() == valuePkgPath {
			continue
		}
		// only recurse into structs within the tree package
		if ft.Kind() == reflect.Struct && isTreePackage(ft.PkgPath()) {
			checkTreeTags(ft, fieldPath, problems, seen)
		}
	}
}
