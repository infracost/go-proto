package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource  `tree:"-"`
	MachineType        value.String                  `tree:"machine_type"`
	InstanceCount      value.Int                     `tree:"instance_count"`
	PurchaseOption     value.Value[PurchaseOption]    `tree:"purchase_option"`
	HasBootDisk        value.Bool                     `tree:"has_boot_disk"`
	BootDiskSizeGB     value.Double                   `tree:"boot_disk_size_gb"`
	BootDiskType       value.Value[DiskType]          `tree:"boot_disk_type"`
	ScratchDisks       value.Int                      `tree:"scratch_disks"`
	GuestAccelerators  []GuestAccelerator             `tree:"guest_accelerators"`
	NATIP              value.String                   `tree:"nat_ip"`
	AttachedDisk       value.String                   `tree:"attached_disk"`
}

type PurchaseOption uint32

const (
	PurchaseOptionUnknown     PurchaseOption = iota
	PurchaseOptionOnDemand
	PurchaseOptionPreemptible
	PurchaseOptionSpot
)

type GuestAccelerator struct {
	Type  value.String `tree:"type"`
	Count value.Int    `tree:"count"`
}
