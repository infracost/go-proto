package compute

import (
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DiskData struct {
	OSType        value.Value[OSType]          `tree:"os_type"`
	DiskType      value.Value[DiskStorageType] `tree:"disk_type"`
	DiskSizeGB    value.Int    `tree:"disk_size_gb"`
	ManagedDiskID value.String `tree:"managed_disk_id"`
}
