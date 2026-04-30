package rds

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource                `tree:"-"`
	Identifier                       value.String                        `tree:"identifier"`
	DBClusterInstanceClass           value.String                        `tree:"db_cluster_instance_class"`
	StorageType                      value.Value[RDSStorageType]         `tree:"storage_type"`
	Engine                           value.Value[RDSEngine]              `tree:"engine"`
	EngineMode                       value.Value[RDSEngineMode]          `tree:"engine_mode"`
	EngineVersion                    value.String                        `tree:"engine_version"`
	BackupRetentionPeriod            value.Int                           `tree:"backup_retention_period"`
	AvailabilityZones                value.List[string]                  `tree:"availability_zones"`
	ServerlessV2ScalingConfiguration []ServerlessV2ScalingConfiguration `tree:"serverless_v2_scaling_configuration"`
	StorageEncrypted                 value.Bool                          `tree:"storage_encrypted"`

	Relationships ClusterRelationships `tree:"-"`
}

type ClusterRelationships struct {
	Instances []*Instance
}

type ServerlessV2ScalingConfiguration struct {
	MinCapacity value.Double `tree:"min_capacity"`
}

type RDSEngineMode uint32

const (
	RDSEngineModeUnknown       RDSEngineMode = iota
	RDSEngineModeProvisioned
	RDSEngineModeServerless
	RDSEngineModeParallelQuery
	RDSEngineModeGlobal
	RDSEngineModeMultiMaster
)
