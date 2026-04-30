package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type WindowsVirtualMachineScaleSet struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	Instances         value.Int    `tree:"instances"`
	LicenseType       value.Value[LicenseType] `tree:"license_type"`
	UltraSSDEnabled   value.Bool   `tree:"ultra_ssd_enabled"`
	OSDisk            *DiskData    `tree:"os_disk"`
}
