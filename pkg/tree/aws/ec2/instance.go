package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource        `tree:"-"`
	Type                     value.String                 `tree:"instance_type"`
	AMI                      value.String                 `tree:"ami"`
	PurchaseOption           value.Value[PurchaseOption]  `tree:"purchase_option"`
	Tenancy                  value.Value[Tenancy]         `tree:"tenancy"`
	EBSOptimized             value.Bool                   `tree:"ebs_optimized"`
	MonitoringEnabled        value.Bool                   `tree:"monitoring_enabled"`
	CPUCredits               value.Value[CPUCreditOption] `tree:"cpu_credits"`
	HostID                   value.String                 `tree:"host_id"`
	LaunchTemplateID         value.String                 `tree:"launch_template_id"`
	LaunchTemplateName       value.String                 `tree:"launch_template_name"`
	AssociatePublicIPAddress value.Bool                   `tree:"associate_public_ip_address"`
	MetadataOptions          InstanceMetadataOptions      `tree:"metadata_options"`
	RootBlockDevice          BlockDeviceMapping           `tree:"root_block_device"`
	// NOTE: Launch templates - which may be linked later (provider-time) - can alter attributes of existing block devices if they share a name
	BlockDeviceMappings []BlockDeviceMapping `tree:"block_device_mappings"`

	Relationships InstanceRelationships `tree:"-"`
}

type InstanceRelationships struct {
	InstanceState  *InstanceStateMapping
	LaunchTemplate *LaunchTemplate
}

type PurchaseOption uint32

const (
	PurchaseOptionUnknown PurchaseOption = iota
	PurchaseOptionOnDemand
	PurchaseOptionSpot
)

type Tenancy uint32

const (
	TenancyUnknown Tenancy = iota
	TenancyShared
	TenancyDedicated
	TenancyHost
)

type CPUCreditOption uint32

const (
	CPUCreditOptionUnknown CPUCreditOption = iota
	CPUCreditOptionStandard
	CPUCreditOptionUnlimited
)

type InstanceMetadataOptions struct {
	HTTPEndpointEnabled value.Bool `tree:"http_endpoint_enabled"`
	HTTPTokensRequired  value.Bool `tree:"http_tokens_required"`
}

type BlockDeviceMapping struct {
	EBSVolume  EBSVolume    `tree:"ebs_volume"`
	DeviceName value.String `tree:"device_name"`
}
