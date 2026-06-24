package s3

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

// CloudFormation/CDK buckets frequently have no explicit BucketName (the physical
// name is generated at deploy time). Their inline configurations carry an empty
// bucket name but share the owning bucket's logical ID, so PostProcess must link
// them to the bucket via that ID fallback rather than the (empty) name.
func TestPostProcess_LinksInlineConfigsForNamelessBucketByID(t *testing.T) {
	const bucketID = "DevArtifactBucket2651DA98"

	s := &S3{
		Buckets: []Bucket{
			{
				Resource: resource.Resource{ID: bucketID},
				Name:     value.New("", 0, "", nil), // no explicit BucketName
			},
		},
		LifecycleConfigurations: []LifecycleConfiguration{
			{
				Resource:   resource.Resource{ID: bucketID}, // inline config shares the bucket's logical ID
				BucketName: value.New("", 0, "", nil),
			},
		},
		IntelligentTieringConfigurations: []IntelligentTieringConfiguration{
			{
				Resource:   resource.Resource{ID: bucketID},
				BucketName: value.New("", 0, "", nil),
			},
		},
		BucketVersioningConfigurations: []BucketVersioningConfiguration{
			{
				Resource:   resource.Resource{ID: bucketID},
				BucketName: value.New("", 0, "", nil),
			},
		},
	}

	s.PostProcess()

	assert.Len(t, s.Buckets[0].Relationships.LifecycleConfigurations, 1,
		"inline lifecycle configuration should link to a nameless bucket via shared ID")
	assert.Len(t, s.Buckets[0].Relationships.IntelligentTieringConfigurations, 1)
	assert.Len(t, s.Buckets[0].Relationships.BucketVersioningConfigurations, 1)
}

// A standalone sub-resource with its own distinct ID and a non-matching bucket
// name must NOT be linked to an unrelated bucket by the ID fallback.
func TestPostProcess_DoesNotFalseMatchStandaloneConfig(t *testing.T) {
	s := &S3{
		Buckets: []Bucket{
			{
				Resource: resource.Resource{ID: "bucket-1"},
				Name:     value.New("bucket-1", 0, "", nil),
			},
		},
		LifecycleConfigurations: []LifecycleConfiguration{
			{
				Resource:   resource.Resource{ID: "standalone-lc"},
				BucketName: value.New("some-other-bucket", 0, "", nil),
			},
		},
	}

	s.PostProcess()

	assert.Empty(t, s.Buckets[0].Relationships.LifecycleConfigurations,
		"a standalone config pointing at another bucket must not link via the ID fallback")
}
