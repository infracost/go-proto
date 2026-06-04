package backup

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Vault struct {
	resource.Resource `tree:"-"`
	KmsKeyArn         value.String `tree:"kms_key_arn"`
}
