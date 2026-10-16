package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Space struct {
	resource.Resource `tree:"-"`
	AppType           value.String `tree:"app_type"`
	EBSVolumeSizeInGB value.Int    `tree:"ebs_volume_size_in_gb"`
}
