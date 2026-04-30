package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EBSSnapshot struct {
	resource.Resource `tree:"-"`
	SizeGB            value.Int    `tree:"size"`
	VolumeID          value.String `tree:"volume_id"`
}
