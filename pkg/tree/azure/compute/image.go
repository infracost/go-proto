package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Image struct {
	resource.Resource      `tree:"-"`
	SourceVirtualMachineID value.String       `tree:"source_virtual_machine_id"`
	OSDiskID               value.String       `tree:"os_disk_id"`
	OSDiskStorageGB        value.Double       `tree:"os_disk_storage_gb"`
	DataDisks              []ImageDataDisk    `tree:"data_disks"`
}

type ImageDataDisk struct {
	ID     value.String `tree:"id"`
	SizeGB value.Double `tree:"size_gb"`
}
