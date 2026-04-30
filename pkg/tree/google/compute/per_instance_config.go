package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type PerInstanceConfig struct {
	resource.Resource        `tree:"-"`
	InstanceGroupManagerRef  value.String `tree:"instance_group_manager_ref"`
}
