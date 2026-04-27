package elasticsearch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Domain struct {
	resource.Resource             `tree:"-"`
	EngineVersion                 value.String `tree:"engine_version"`
	ClusterInstanceType           value.String `tree:"cluster_instance_type"`
	ClusterInstanceCount          value.Int    `tree:"cluster_instance_count"`
	ClusterDedicatedMasterCount   value.Int    `tree:"cluster_dedicated_master_count"`
	ClusterDedicatedMasterType    value.String `tree:"cluster_dedicated_master_type"`
	ClusterDedicatedMasterEnabled value.Bool   `tree:"cluster_dedicated_master_enabled"`
	ClusterWarmEnabled            value.Bool   `tree:"cluster_warm_enabled"`
	ClusterWarmType               value.String `tree:"cluster_warm_type"`
	ClusterWarmCount              value.Int    `tree:"cluster_warm_count"`
	EBSEnabled                    value.Bool   `tree:"ebs_enabled"`
	EBSVolumeType                 value.Value[EBSVolumeType] `tree:"ebs_volume_type"`
	EBSVolumeSize                 value.Int    `tree:"ebs_volume_size"`
	EBSVolumeIOPS                 value.Int    `tree:"ebs_volume_iops"`
	EBSVolumeThroughput           value.Int    `tree:"ebs_volume_throughput"`
	IsOpenSearch                  value.Bool   `tree:"is_opensearch"`
}

type EBSVolumeType uint32

const (
	EBSVolumeTypeGP2      EBSVolumeType = 0
	EBSVolumeTypeGP3      EBSVolumeType = 1
	EBSVolumeTypeIO1      EBSVolumeType = 2
	EBSVolumeTypeStandard EBSVolumeType = 3
)
