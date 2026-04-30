package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LinuxVirtualMachineScaleSet struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	Instances         value.Int    `tree:"instances"`
	UltraSSDEnabled   value.Bool   `tree:"ultra_ssd_enabled"`
	OSDisk            *DiskData    `tree:"os_disk"`
}
