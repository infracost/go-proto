package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EBSSnapshotImport struct {
	resource.Resource `tree:"-"`
	StorageTier       value.String `tree:"storage_tier"`
	SizeGB            value.Double `tree:"size"`
}
