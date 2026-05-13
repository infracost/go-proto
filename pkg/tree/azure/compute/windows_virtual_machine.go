package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type WindowsVirtualMachine struct {
	resource.Resource `tree:"-"`
	Size              value.String `tree:"size"`
	LicenseType       value.Value[LicenseType] `tree:"license_type"`
	Zone              value.String `tree:"zone"`
	EncryptionAtHost  value.Bool   `tree:"encryption_at_host"`
	UltraSSDEnabled   value.Bool   `tree:"ultra_ssd_enabled"`
	AvailabilitySetID value.String `tree:"availability_set_id"`
	OSDisk            *DiskData    `tree:"os_disk"`
}
