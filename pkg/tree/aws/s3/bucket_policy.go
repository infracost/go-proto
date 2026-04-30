package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketPolicy struct {
	resource.Resource       `tree:"-"`
	BucketName              value.String `tree:"bucket"`
	DeniesInsecureTransport value.Bool   `tree:"denies_insecure_transport"`
}
