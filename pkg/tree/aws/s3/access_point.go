package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type AccessPoint struct {
	resource.Resource              `tree:"-"`
	Name                           value.String                    `tree:"name"`
	Bucket                         value.String                    `tree:"bucket"`
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `tree:"public_access_block_configuration"`
}

type PublicAccessBlockConfiguration struct {
	BlockPublicACLs       value.Bool `tree:"block_public_acls"`
	BlockPublicPolicy     value.Bool `tree:"block_public_policy"`
	IgnorePublicACLs      value.Bool `tree:"ignore_public_acls"`
	RestrictPublicBuckets value.Bool `tree:"restrict_public_buckets"`
}
