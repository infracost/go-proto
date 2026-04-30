package emr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InstanceFleet struct {
	resource.Resource   `tree:"-"`
	ClusterID           value.String         `tree:"cluster_id"`
	InstanceTypeConfigs []InstanceTypeConfig `tree:"instance_type_configs"`
}

type InstanceTypeConfig struct {
	InstanceType value.String `tree:"instance_type"`
	EBSConfigs   []EBSConfig  `tree:"ebs_config"`
}
