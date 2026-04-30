package emr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InstanceGroup struct {
	resource.Resource `tree:"-"`
	ClusterID         value.String `tree:"cluster_id"`
	InstanceType      value.String `tree:"instance_type"`
	InstanceCount     value.Int    `tree:"instance_count"`
	Name              value.String `tree:"name"`
	EBSConfigs        []EBSConfig  `tree:"ebs_config"`
}
