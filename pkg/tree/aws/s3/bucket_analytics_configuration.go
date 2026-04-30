package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketAnalyticsConfiguration struct {
	resource.Resource `tree:"-"`
	Bucket            value.String `tree:"bucket"`
	BucketName        value.String `tree:"bucket_name"`
	Name              value.String `tree:"name"`
}
