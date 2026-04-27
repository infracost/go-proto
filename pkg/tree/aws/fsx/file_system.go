package fsx

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FileSystem struct {
	resource.Resource   `tree:"-"`
	Type                value.Value[FileSystemType] `tree:"type"`
	StorageType         value.Value[StorageType]    `tree:"storage_type"`
	ThroughputCapacity  value.Int                   `tree:"throughput_capacity"`
	ProvisionedIOPS     value.Int                   `tree:"provisioned_iops"`
	ProvisionedIOPSMode value.String                `tree:"provisioned_iops_mode"`
	StorageCapacityGB   value.Int                   `tree:"storage_capacity_gb"`
	DeploymentType      value.Value[DeploymentType] `tree:"deployment_type"`
	DataCompression     value.Value[DataCompression] `tree:"data_compression"`
}

type FileSystemType uint32

const (
	FileSystemTypeWindows FileSystemType = 0
	FileSystemTypeOpenZFS FileSystemType = 1
	FileSystemTypeLustre  FileSystemType = 2
	FileSystemTypeOntap   FileSystemType = 3
)

type StorageType uint32

const (
	StorageTypeSSD                StorageType = 0
	StorageTypeHDD                StorageType = 1
	StorageTypeIntelligentTiering StorageType = 2
)

type DeploymentType uint32

const (
	DeploymentTypeSingleAZ1   DeploymentType = 0
	DeploymentTypeSingleAZ2   DeploymentType = 1
	DeploymentTypeSingleAZHA1 DeploymentType = 2
	DeploymentTypeSingleAZHA2 DeploymentType = 3
	DeploymentTypeMultiAZ1    DeploymentType = 4
)

type DataCompression uint32

const (
	DataCompressionNone DataCompression = 0
	DataCompressionZSTD DataCompression = 1
)
