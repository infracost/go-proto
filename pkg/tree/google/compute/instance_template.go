package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InstanceTemplate struct {
	resource.Resource  `tree:"-"`
	Name               value.String       `tree:"name"`
	MachineType        value.String       `tree:"machine_type"`
	ScratchDisks       value.Int          `tree:"scratch_disks"`
	DiskData           []DiskData         `tree:"disk_data"`
	GuestAccelerators  []GuestAccelerator `tree:"guest_accelerators"`
}

type DiskData struct {
	Type      value.Value[DiskType] `tree:"type"`
	StorageGB value.Double          `tree:"storage_gb"`
}
