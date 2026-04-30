package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VirtualMachine struct {
	resource.Resource              `tree:"-"`
	StorageImageReferenceOffer     value.String `tree:"storage_image_reference_offer"`
	VMSize                         value.String `tree:"vm_size"`
	LicenseType                    value.Value[LicenseType] `tree:"license_type"`
	Zone                           value.String `tree:"zone"`
	EncryptionAtHost               value.Bool   `tree:"encryption_at_host"`
	OSDisk                         *DiskData    `tree:"os_disk"`
	StorageOSDisk                  *DiskData    `tree:"storage_os_disk"`
	StorageDataDisks               []*DiskData  `tree:"storage_data_disks"`
}
