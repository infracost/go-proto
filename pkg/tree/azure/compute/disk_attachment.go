package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DiskAttachment struct {
	resource.Resource `tree:"-"`
	ManagedDiskID     value.String `tree:"managed_disk_id"`
	VirtualMachineID  value.String `tree:"virtual_machine_id"`
}
