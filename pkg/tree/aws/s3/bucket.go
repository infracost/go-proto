package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Bucket struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	ObjectTagsEnabled value.Bool   `tree:"object_tags_enabled"`

	Relationships BucketRelationships `tree:"-"`
}

type BucketRelationships struct {
	LifecycleConfigurations          []*LifecycleConfiguration
	IntelligentTieringConfigurations []*IntelligentTieringConfiguration
	BucketVersioningConfigurations   []*BucketVersioningConfiguration
	BucketPolicies                   []*BucketPolicy
}
