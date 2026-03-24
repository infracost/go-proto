package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LaunchTemplate struct {
	resource.Resource                 `tree:"-"`
	Name                              value.String                 `tree:"name"`
	InstanceType                      value.String                 `tree:"instance_type"`
	AMI                               value.String                 `tree:"ami"`
	EBSOptimized                      value.Bool                   `tree:"ebs_optimised"`
	MonitoringEnabled                 value.Bool                   `tree:"monitoring_enabled"`
	CPUCredits                        value.Value[CPUCreditOption] `tree:"cpu_credits"`
	Tenancy                           value.Value[Tenancy]         `tree:"tenancy"`
	MarketType                        value.Value[PurchaseOption]  `tree:"market_type"`
	BlockDeviceMappings               []BlockDeviceMapping         `tree:"block_device_mappings"`
	NetworkInterfaces                 []NetworkInterface           `tree:"network_interfaces"`
	OnDemandPercentageAboveBaseCount  value.Int                    `tree:"on_demand_percentage_above_base_count"`
	MetadataOptions                   InstanceMetadataOptions      `tree:"metadata_options"`
	InstanceTagSpecifications         resource.Tags                `tree:"instance_tag_specifications"`
	VolumeTagSpecifications           resource.Tags                `tree:"volume_tag_specifications"`
	NetworkInterfaceTagSpecifications resource.Tags                `tree:"network_interface_tag_specifications"`
}

type NetworkInterface struct {
	AssociatePublicIPAddress value.Bool
	DeviceIndex              value.Int
}
