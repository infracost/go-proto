package rds

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource                  `tree:"-"`
	Identifier                         value.String                `tree:"identifier"`
	ReplicateSourceDBIdentifier        value.String                `tree:"replicate_source_db"`
	StorageType                        value.Value[RDSStorageType] `tree:"storage_type"`
	LicenseModel                       value.Value[RDSLicenseModel] `tree:"license_model"`
	InstanceClass                      value.String                 `tree:"instance_class"`
	Engine                             value.Value[RDSEngine]       `tree:"engine"`
	EngineVersion                      value.String                `tree:"engine_version"`
	ClusterID                          value.String                `tree:"cluster_identifier"`
	IOPS                               value.Int                   `tree:"iops"`
	MultiAZ                            value.Bool                  `tree:"multi_az"`
	PerformanceInsightsEnabled         value.Bool                  `tree:"performance_insights_enabled"`
	PerformanceInsightsRetentionPeriod value.Int                   `tree:"performance_insights_retention_period"`
	IOOptimized                        value.Bool                  `tree:"io_optimized"`
	BackupRetentionPeriod              value.Int                   `tree:"backup_retention_period"`
	AllocatedStorageGB                 value.Double                `tree:"allocated_storage"`
	StorageEncrypted                   value.Bool                  `tree:"storage_encrypted"`
	IsClusterInstance                  value.Bool                  `tree:"is_cluster_instance"`

	Relationships InstanceRelationships `tree:"-"`
}

type InstanceRelationships struct {
	ReplicateSourceDB *Instance
	Cluster           *Cluster
}

type RDSStorageType uint32

const (
	RDSStorageTypeUnknown RDSStorageType = iota
	RDSStorageTypeGP2
	RDSStorageTypeGP3
	RDSStorageTypeIO1
	RDSStorageTypeIO2
	RDSStorageTypeStandard
	RDSStorageTypeAurora
	RDSStorageTypeAuroraIOPT1
)
