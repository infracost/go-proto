package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TaskDefinition struct {
	resource.Resource       `tree:"-"`
	Family                  value.String   `tree:"family"`
	Memory                  value.String   `tree:"memory"`
	CPU                     value.String   `tree:"cpu"`
	CPUArchitecture         value.String   `tree:"cpu_architecture"`
	RequiresCompatibilities []value.String `tree:"requires_compatibilities"`
	InferenceAccelerators   []InferenceAccelerator `tree:"inference_accelerators"`
}

type InferenceAccelerator struct {
	DeviceType value.Value[InferenceAcceleratorDeviceType] `tree:"device_type"`
	DeviceName value.String                                `tree:"device_name"`
}

type InferenceAcceleratorDeviceType uint32

const (
	InferenceAcceleratorDeviceTypeEIA1Medium InferenceAcceleratorDeviceType = 0
	InferenceAcceleratorDeviceTypeEIA1Large  InferenceAcceleratorDeviceType = 1
	InferenceAcceleratorDeviceTypeEIA1XLarge InferenceAcceleratorDeviceType = 2
	InferenceAcceleratorDeviceTypeEIA2Medium InferenceAcceleratorDeviceType = 3
	InferenceAcceleratorDeviceTypeEIA2Large  InferenceAcceleratorDeviceType = 4
	InferenceAcceleratorDeviceTypeEIA2XLarge InferenceAcceleratorDeviceType = 5
)
