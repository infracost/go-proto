package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TaskDefinition struct {
	resource.Resource       `tree:"-"`
	Family                  value.String                 `tree:"family"`
	MemoryMiB               value.Int                    `tree:"memory"`
	CPU                     value.Int                    `tree:"cpu"`
	CPUArchitecture         value.Value[CPUArchitecture] `tree:"cpu_architecture"`
	RequiresCompatibilities value.List[Compatibility]    `tree:"requires_compatibilities"`
	InferenceAccelerators   []InferenceAccelerator       `tree:"inference_accelerators"`
	ContainerDefinitions    []ContainerDefinition        `tree:"container_definitions"`

	Relationships TaskDefinitionRelationships `tree:"-"`
}

type ContainerDefinition struct {
	Name                 value.String          `tree:"name"`
	Image                value.String          `tree:"image"`
	EnvironmentVariables []EnvironmentVariable `tree:"environment_variables"`
}

type EnvironmentVariable struct {
	Name  value.String `tree:"name"`
	Value value.String `tree:"value"`
}

type TaskDefinitionRelationships struct{}

type InferenceAccelerator struct {
	DeviceType value.Value[InferenceAcceleratorDeviceType] `tree:"device_type"`
	DeviceName value.String                                `tree:"device_name"`
}

type InferenceAcceleratorDeviceType uint32

const (
	InferenceAcceleratorDeviceTypeUnknown InferenceAcceleratorDeviceType = iota
	InferenceAcceleratorDeviceTypeEIA1Medium
	InferenceAcceleratorDeviceTypeEIA1Large
	InferenceAcceleratorDeviceTypeEIA1XLarge
	InferenceAcceleratorDeviceTypeEIA2Medium
	InferenceAcceleratorDeviceTypeEIA2Large
	InferenceAcceleratorDeviceTypeEIA2XLarge
)

type CPUArchitecture uint32

const (
	CPUArchitectureUnknown CPUArchitecture = iota
	CPUArchitectureX86_64
	CPUArchitectureARM64
)

type Compatibility uint32

const (
	CompatibilityUnknown Compatibility = iota
	CompatibilityEC2
	CompatibilityFargate
)
