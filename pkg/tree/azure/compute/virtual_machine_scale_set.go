package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VirtualMachineScaleSet struct {
	resource.Resource `tree:"-"`
	SKUName           value.String             `tree:"sku_name"`
	SKUCapacity       value.Int                `tree:"sku_capacity"`
	Instances         value.Int                `tree:"instances"`
	LicenseType       value.Value[LicenseType] `tree:"license_type"`
	IsWindows         value.Bool               `tree:"is_windows"`
	StorageOSDisk     *DiskData                `tree:"storage_os_disk"`
	StorageDataDisks  []*DiskData              `tree:"storage_data_disks"`
}
