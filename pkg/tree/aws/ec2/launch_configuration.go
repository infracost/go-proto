package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LaunchConfiguration struct {
	resource.Resource   `tree:"-"`
	Name                value.String                 `tree:"name"`
	Tenancy             value.Value[Tenancy]         `tree:"tenancy"`
	PurchaseOption      value.Value[PurchaseOption]  `tree:"purchase_option"`
	AMI                 value.String                 `tree:"ami"`
	InstanceType        value.String                 `tree:"instance_type"`
	MonitoringEnabled   value.Bool                   `tree:"monitoring_enabled"`
	EBSOptimized        value.Bool                   `tree:"ebs_optimized"`
	CPUCredits          value.Value[CPUCreditOption] `tree:"cpu_credits"`
	RootBlockDevice     BlockDeviceMapping           `tree:"root_block_device"`
	BlockDeviceMappings []BlockDeviceMapping         `tree:"block_device_mappings"`
}
