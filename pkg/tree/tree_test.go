package tree

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/aws/rds"
	"github.com/infracost/go-proto/pkg/tree/aws/s3"
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

// TestPostProcess_NoDoubleInvocation verifies that calling Tree.PostProcess()
// does not invoke service-level PostProcess methods twice. The reflective walker
// in tree.go discovers and calls PostProcess on each service struct. Provider-level
// PostProcess methods (AWS, Azure, Google) must NOT also call service-level
// PostProcess, or relationships will be duplicated.
func TestPostProcess_NoDoubleInvocation(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			S3: s3.S3{
				Buckets: []s3.Bucket{
					{Resource: resource.Resource{ID: "bucket-1"}, Name: value.New("my-bucket", 0, "", nil)},
				},
				BucketVersioningConfigurations: []s3.BucketVersioningConfiguration{
					{Resource: resource.Resource{ID: "vc-1"}, BucketName: value.New("my-bucket", 0, "", nil), Enabled: value.New(true, 0, "", nil)},
				},
				LifecycleConfigurations: []s3.LifecycleConfiguration{
					{Resource: resource.Resource{ID: "lc-1"}, BucketName: value.New("my-bucket", 0, "", nil)},
				},
				BucketPolicies: []s3.BucketPolicy{
					{Resource: resource.Resource{ID: "bp-1"}, BucketName: value.New("my-bucket", 0, "", nil)},
				},
				IntelligentTieringConfigurations: []s3.IntelligentTieringConfiguration{
					{Resource: resource.Resource{ID: "it-1"}, BucketName: value.New("my-bucket", 0, "", nil)},
				},
			},
			RDS: rds.RDS{
				Clusters: []rds.Cluster{
					{Resource: resource.Resource{ID: "cluster-1"}, Identifier: value.New("my-cluster", 0, "", nil)},
				},
				Instances: []rds.Instance{
					{Resource: resource.Resource{ID: "inst-1"}, ClusterID: value.New("my-cluster", 0, "", nil)},
				},
			},
		},
	}

	tree.PostProcess()

	// S3: each relationship should be linked exactly once
	bucket := tree.AWS.S3.Buckets[0]
	assert.Len(t, bucket.Relationships.BucketVersioningConfigurations, 1, "versioning should be linked exactly once")
	assert.Len(t, bucket.Relationships.LifecycleConfigurations, 1, "lifecycle should be linked exactly once")
	assert.Len(t, bucket.Relationships.BucketPolicies, 1, "policy should be linked exactly once")
	assert.Len(t, bucket.Relationships.IntelligentTieringConfigurations, 1, "intelligent tiering should be linked exactly once")

	// RDS: instance→cluster should be linked exactly once
	inst := tree.AWS.RDS.Instances[0]
	require.NotNil(t, inst.Relationships.Cluster, "RDS instance should be linked to cluster")
	assert.Len(t, tree.AWS.RDS.Clusters[0].Relationships.Instances, 1, "cluster should have exactly one instance")
}

// TestPostProcess_MultipleCallsDoNotPanic verifies that calling PostProcess
// multiple times does not panic.
func TestPostProcess_MultipleCallsDoNotPanic(t *testing.T) {
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

	// Should not panic on multiple calls
	tree.PostProcess()
	tree.PostProcess()
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
