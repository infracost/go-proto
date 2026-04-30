package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketVersioningConfiguration struct {
	resource.Resource `tree:"-"`
	BucketName        value.String `tree:"bucket"`
	Enabled           value.Bool   `tree:"enabled"`
}
