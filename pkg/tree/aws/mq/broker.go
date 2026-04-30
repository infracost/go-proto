package mq

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Broker struct {
	resource.Resource `tree:"-"`
	EngineType        value.Value[EngineType]     `tree:"engine_type"`
	HostInstanceType  value.String                `tree:"host_instance_type"`
	StorageType       value.Value[StorageType]    `tree:"storage_type"`
	DeploymentMode    value.Value[DeploymentMode] `tree:"deployment_mode"`
}

type EngineType uint32

const (
	EngineTypeUnknown  EngineType = iota
	EngineTypeActiveMQ
	EngineTypeRabbitMQ
)

type StorageType uint32

const (
	StorageTypeUnknown StorageType = iota
	StorageTypeEFS
	StorageTypeEBS
)

type DeploymentMode uint32

const (
	DeploymentModeUnknown                DeploymentMode = iota
	DeploymentModeSingleInstance
	DeploymentModeActiveStandbyMultiAZ
	DeploymentModeClusterMultiAZ
)
