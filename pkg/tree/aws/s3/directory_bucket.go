package s3

import "github.com/infracost/go-proto/pkg/tree/resource"

// DirectoryBucket is an S3 Express One Zone bucket.
type DirectoryBucket struct {
	resource.Resource `tree:"-"`
}
