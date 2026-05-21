package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InstanceGroupManager struct {
	resource.Resource  `tree:"-"`
	Name               value.String                `tree:"name"`
	MachineType        value.String                `tree:"machine_type"`
	PurchaseOption     value.Value[PurchaseOption]  `tree:"purchase_option"`
	TargetSize         value.Int                   `tree:"target_size"`
	ScratchDisks       value.Int                   `tree:"scratch_disks"`
	GuestAccelerators  []GuestAccelerator          `tree:"guest_accelerators"`
	DiskData           []DiskData                  `tree:"disk_data"`
	TemplateRef        value.String                `tree:"template_ref"`

	Relationships InstanceGroupManagerRelationships `tree:"-"`
}

type InstanceGroupManagerRelationships struct {
	InstanceTemplate   *InstanceTemplate
	PerInstanceConfigs []*PerInstanceConfig
}
