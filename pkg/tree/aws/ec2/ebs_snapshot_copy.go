package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EBSSnapshotCopy struct {
	resource.Resource  `tree:"-"`
	SizeGB             value.Int    `tree:"size"`
	SourceSnapshotID   value.String `tree:"source_snapshot_id"`
	SourceRegion       value.String `tree:"source_region"`
}
