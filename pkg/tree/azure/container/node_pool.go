package container

import (
	"github.com/infracost/go-proto/pkg/tree/value"
)

type NodePool struct {
	VMSize       value.String `tree:"vm_size"`
	NodeCount    value.Int    `tree:"node_count"`
	MinCount     value.Int    `tree:"min_count"`
	MaxCount     value.Int    `tree:"max_count"`
	OSType       value.Value[OSType]     `tree:"os_type"`
	OSDiskType   value.Value[OSDiskType] `tree:"os_disk_type"`
	OSDiskSizeGB value.Int    `tree:"os_disk_size_gb"`
}
