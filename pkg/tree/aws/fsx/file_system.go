package fsx

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FileSystem struct {
	resource.Resource        `tree:"-"`
	Type                     value.Value[FileSystemType]      `tree:"type"`
	StorageType              value.Value[StorageType]         `tree:"storage_type"`
	ThroughputCapacityMBperS value.Int                        `tree:"throughput_capacity"`
	ProvisionedIOPS          value.Int                        `tree:"provisioned_iops"`
	ProvisionedIOPSMode      value.Value[ProvisionedIOPSMode] `tree:"provisioned_iops_mode"`
	StorageCapacityGB        value.Int                        `tree:"storage_capacity_gb"`
	DeploymentType           value.Value[DeploymentType]      `tree:"deployment_type"`
	DataCompression          value.Value[DataCompression]     `tree:"data_compression"`
}

type FileSystemType uint32

const (
	FileSystemTypeUnknown FileSystemType = iota
	FileSystemTypeWindows
	FileSystemTypeOpenZFS
	FileSystemTypeLustre
	FileSystemTypeOntap
)

type StorageType uint32

const (
	StorageTypeUnknown StorageType = iota
	StorageTypeSSD
	StorageTypeHDD
	StorageTypeIntelligentTiering
)

type DeploymentType uint32

const (
	DeploymentTypeUnknown DeploymentType = iota
	DeploymentTypeSingleAZ1
	DeploymentTypeSingleAZ2
	DeploymentTypeSingleAZHA1
	DeploymentTypeSingleAZHA2
	DeploymentTypeMultiAZ1
)

type DataCompression uint32

const (
	DataCompressionUnknown DataCompression = iota
	DataCompressionNone
	DataCompressionZSTD
)

type ProvisionedIOPSMode uint32

const (
	ProvisionedIOPSModeUnknown ProvisionedIOPSMode = iota
	ProvisionedIOPSModeAutomatic
	ProvisionedIOPSModeUserProvisioned
)
