package eks

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type NodeGroup struct {
	resource.Resource  `tree:"-"`
	InstanceCount      value.Int      `tree:"instance_count"`
	DiskSize           value.Int      `tree:"disk_size"`
	Name               value.String   `tree:"name"`
	ClusterName        value.String   `tree:"cluster_name"`
	InstanceTypes      []value.String `tree:"instance_types"`
	PurchaseOption     value.Value[PurchaseOption] `tree:"purchase_option"`
	LaunchTemplateID   value.String   `tree:"launch_template_id"`
	LaunchTemplateName value.String   `tree:"launch_template_name"`
}

type PurchaseOption uint32

const (
	PurchaseOptionOnDemand PurchaseOption = 0
	PurchaseOptionSpot     PurchaseOption = 1
)
