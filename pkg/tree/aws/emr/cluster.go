package emr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type Cluster struct {
	resource.Resource   `tree:"-"`
	MasterInstanceGroup *InstanceGroup `tree:"master_instance_group"`
	MasterInstanceFleet *InstanceFleet `tree:"master_instance_fleet"`
	CoreInstanceGroup   *InstanceGroup `tree:"core_instance_group"`
	CoreInstanceFleet   *InstanceFleet `tree:"core_instance_fleet"`
	TaskInstanceGroup   *InstanceGroup `tree:"task_instance_group"`
	TaskInstanceFleet   *InstanceFleet `tree:"task_instance_fleet"`
}
