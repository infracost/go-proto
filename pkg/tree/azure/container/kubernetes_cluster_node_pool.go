package container

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type KubernetesClusterNodePool struct {
	resource.Resource `tree:"-"`
	VMSize            value.String `tree:"vm_size"`
	NodeCount         value.Int    `tree:"node_count"`
	MinCount          value.Int    `tree:"min_count"`
	MaxCount          value.Int    `tree:"max_count"`
	OSType            value.Value[OSType]     `tree:"os_type"`
	OSDiskType        value.Value[OSDiskType] `tree:"os_disk_type"`
	OSDiskSizeGB      value.Int    `tree:"os_disk_size_gb"`
	ClusterID         value.String `tree:"cluster_id"`

	Relationships KubernetesClusterNodePoolRelationships `tree:"-"`
}

type KubernetesClusterNodePoolRelationships struct {
	Cluster *KubernetesCluster
}
