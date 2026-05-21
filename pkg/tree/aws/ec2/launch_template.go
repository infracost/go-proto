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
	OnDemandPercentageAboveBaseCount  value.Int                                    `tree:"on_demand_percentage_above_base_count"`
	ElasticInferenceAcceleratorType   value.Value[ElasticInferenceAcceleratorType] `tree:"elastic_inference_accelerator_type"`
	MetadataOptions                   InstanceMetadataOptions                      `tree:"metadata_options"`
	InstanceTagSpecifications         resource.Tags                `tree:"instance_tag_specifications"`
	VolumeTagSpecifications           resource.Tags                `tree:"volume_tag_specifications"`
	NetworkInterfaceTagSpecifications resource.Tags                `tree:"network_interface_tag_specifications"`
}

type NetworkInterface struct {
	AssociatePublicIPAddress value.Bool `tree:"associate_public_ip_address"`
	DeviceIndex              value.Int  `tree:"device_index"`
}

// ElasticInferenceAcceleratorType is the closed set of Elastic Inference
// accelerator device types. AWS deprecated Elastic Inference in April 2024,
// so no new variants are expected.
type ElasticInferenceAcceleratorType uint32

const (
	ElasticInferenceAcceleratorTypeUnknown  ElasticInferenceAcceleratorType = iota
	ElasticInferenceAcceleratorTypeEIA1Medium
	ElasticInferenceAcceleratorTypeEIA1Large
	ElasticInferenceAcceleratorTypeEIA1XLarge
	ElasticInferenceAcceleratorTypeEIA2Medium
	ElasticInferenceAcceleratorTypeEIA2Large
	ElasticInferenceAcceleratorTypeEIA2XLarge
)

// String returns the AWS-canonical device type name (e.g. "eia2.medium").
func (t ElasticInferenceAcceleratorType) String() string {
	switch t {
	case ElasticInferenceAcceleratorTypeEIA1Medium:
		return "eia1.medium"
	case ElasticInferenceAcceleratorTypeEIA1Large:
		return "eia1.large"
	case ElasticInferenceAcceleratorTypeEIA1XLarge:
		return "eia1.xlarge"
	case ElasticInferenceAcceleratorTypeEIA2Medium:
		return "eia2.medium"
	case ElasticInferenceAcceleratorTypeEIA2Large:
		return "eia2.large"
	case ElasticInferenceAcceleratorTypeEIA2XLarge:
		return "eia2.xlarge"
	}
	return ""
}
