package elasticsearch

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Domain struct {
	resource.Resource             `tree:"-"`
	EngineVersion                 value.String                   `tree:"engine_version"`
	ClusterInstanceType           value.String                   `tree:"cluster_instance_type"`
	ClusterInstanceCount          value.Int                      `tree:"cluster_instance_count"`
	ClusterDedicatedMasterCount   value.Int                      `tree:"cluster_dedicated_master_count"`
	ClusterDedicatedMasterType    value.String                   `tree:"cluster_dedicated_master_type"`
	ClusterDedicatedMasterEnabled value.Bool                     `tree:"cluster_dedicated_master_enabled"`
	ClusterWarmEnabled            value.Bool                     `tree:"cluster_warm_enabled"`
	ClusterWarmType               value.String                   `tree:"cluster_warm_type"`
	ClusterWarmCount              value.Int                      `tree:"cluster_warm_count"`
	EBSEnabled                    value.Bool                     `tree:"ebs_enabled"`
	EBSVolumeType                 value.Value[ec2.EBSVolumeType] `tree:"ebs_volume_type"`
	EBSVolumeSizeGB               value.Int                      `tree:"ebs_volume_size"`
	EBSVolumeIOPS                 value.Int                      `tree:"ebs_volume_iops"`
	EBSVolumeThroughputMiBperS    value.Int                      `tree:"ebs_volume_throughput"`
	IsOpenSearch                  value.Bool                     `tree:"is_opensearch"`
	EncryptAtRestEnabled          value.Bool                     `tree:"encrypt_at_rest_enabled"`
	NodeToNodeEncryptionEnabled   value.Bool                     `tree:"node_to_node_encryption_enabled"`
	EnforceHTTPS                  value.Bool                     `tree:"enforce_https"`
	TLSSecurityPolicy             value.String                   `tree:"tls_security_policy"`
}

type EBSVolumeType uint32

const (
	EBSVolumeTypeUnknown EBSVolumeType = iota
	EBSVolumeTypeGP2
	EBSVolumeTypeGP3
	EBSVolumeTypeIO1
	EBSVolumeTypeStandard
)
