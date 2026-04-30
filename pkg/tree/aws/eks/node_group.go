package eks

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type NodeGroup struct {
	resource.Resource  `tree:"-"`
	InstanceCount      value.Int                   `tree:"instance_count"`
	DiskSizeGB         value.Int                   `tree:"disk_size"`
	Name               value.String                `tree:"name"`
	ClusterName        value.String                `tree:"cluster_name"`
	InstanceTypes      value.List[string]          `tree:"instance_types"`
	PurchaseOption     value.Value[PurchaseOption] `tree:"purchase_option"`
	LaunchTemplateID   value.String                `tree:"launch_template_id"`
	LaunchTemplateName value.String                `tree:"launch_template_name"`

	Relationships NodeGroupRelationships `tree:"-"`
}

type NodeGroupRelationships struct {
	LaunchTemplate *ec2.LaunchTemplate `tree:"launch_template"`
}

type PurchaseOption uint32

const (
	PurchaseOptionUnknown PurchaseOption = iota
	PurchaseOptionOnDemand
	PurchaseOptionSpot
)
