package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ManagedDisk struct {
	resource.Resource `tree:"-"`
	DiskType          value.Value[DiskStorageType] `tree:"disk_type"`
	DiskSizeGB        value.Int    `tree:"disk_size_gb"`
	IOPS              value.Int    `tree:"iops"`
	MBpsReadWrite     value.Int    `tree:"mbps_read_write"`

	Relationships ManagedDiskRelationships `tree:"-"`
}

// ManagedDiskRelationships links a managed disk to whichever VM resource it
// is attached to (only one of the three pointers should be non-nil).
type ManagedDiskRelationships struct {
	VirtualMachine        *VirtualMachine
	LinuxVirtualMachine   *LinuxVirtualMachine
	WindowsVirtualMachine *WindowsVirtualMachine
}
