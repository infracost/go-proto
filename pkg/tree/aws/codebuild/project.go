package codebuild

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Project struct {
	resource.Resource `tree:"-"`
	ComputeType       value.Value[ComputeType]     `tree:"compute_type"`
	EnvironmentType   value.Value[EnvironmentType] `tree:"environment_type"`
}

type ComputeType uint32

const (
	ComputeTypeUnknown ComputeType = iota
	ComputeTypeGeneral1Small
	ComputeTypeGeneral1Medium
	ComputeTypeGeneral1Large
	ComputeTypeGeneral12XL
	ComputeTypeArmLambda1GB
	ComputeTypeArmLambda2GB
	ComputeTypeArmLambda4GB
	ComputeTypeArmLambda8GB
	ComputeTypeArmLambda10GB
	ComputeTypeLinuxLambda1GB
	ComputeTypeLinuxLambda2GB
	ComputeTypeLinuxLambda4GB
	ComputeTypeLinuxLambda8GB
	ComputeTypeLinuxLambda10GB
)

type EnvironmentType uint32

const (
	EnvironmentTypeUnknown EnvironmentType = iota
	EnvironmentTypeLinuxContainer
	EnvironmentTypeLinuxGPUContainer
	EnvironmentTypeARMContainer
	EnvironmentTypeWindowsServer2019Container
	EnvironmentTypeArmEC2Container
	EnvironmentTypeLinuxEC2Container
	EnvironmentTypeWindowsServer2022Container
	EnvironmentTypeWindowsEC2Container
	EnvironmentTypeArmLambdaContainer
	EnvironmentTypeLinuxLambdaContainer
)
