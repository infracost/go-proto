package emr

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EBSConfig struct {
	Type               value.Value[ec2.EBSVolumeType] `tree:"type"`
	IOPS               value.Int    `tree:"iops"`
	SizeGB             value.Int    `tree:"size"`
	VolumesPerInstance value.Int    `tree:"volumes_per_instance"`
}
