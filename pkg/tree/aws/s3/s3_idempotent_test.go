package s3

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &S3{
		Buckets: []Bucket{
			{
				Resource: resource.Resource{ID: "bucket-1"},
				Name:     value.New("bucket-1", 0, "", nil),
			},
		},
		LifecycleConfigurations: []LifecycleConfiguration{
			{
				Resource:   resource.Resource{ID: "lc-1"},
				BucketName: value.New("bucket-1", 0, "", nil),
			},
		},
		IntelligentTieringConfigurations: []IntelligentTieringConfiguration{
			{
				Resource:   resource.Resource{ID: "itc-1"},
				BucketName: value.New("bucket-1", 0, "", nil),
			},
		},
		BucketVersioningConfigurations: []BucketVersioningConfiguration{
			{
				Resource:   resource.Resource{ID: "bvc-1"},
				BucketName: value.New("bucket-1", 0, "", nil),
			},
		},
		BucketPolicies: []BucketPolicy{
			{
				Resource:   resource.Resource{ID: "bp-1"},
				BucketName: value.New("bucket-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	lcs := append([]*LifecycleConfiguration(nil), s.Buckets[0].Relationships.LifecycleConfigurations...)
	itcs := append([]*IntelligentTieringConfiguration(nil), s.Buckets[0].Relationships.IntelligentTieringConfigurations...)
	bvcs := append([]*BucketVersioningConfiguration(nil), s.Buckets[0].Relationships.BucketVersioningConfigurations...)
	bps := append([]*BucketPolicy(nil), s.Buckets[0].Relationships.BucketPolicies...)

	s.PostProcess()
	assert.Equal(t, lcs, s.Buckets[0].Relationships.LifecycleConfigurations)
	assert.Equal(t, itcs, s.Buckets[0].Relationships.IntelligentTieringConfigurations)
	assert.Equal(t, bvcs, s.Buckets[0].Relationships.BucketVersioningConfigurations)
	assert.Equal(t, bps, s.Buckets[0].Relationships.BucketPolicies)

	assert.Len(t, lcs, 1)
	assert.Len(t, itcs, 1)
	assert.Len(t, bvcs, 1)
	assert.Len(t, bps, 1)
}
