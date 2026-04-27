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
	ComputeTypeUnknown        ComputeType = 0
	ComputeTypeGeneral1Small  ComputeType = 1
	ComputeTypeGeneral1Medium ComputeType = 2
	ComputeTypeGeneral1Large  ComputeType = 3
	ComputeTypeGeneral12XL    ComputeType = 4
	ComputeTypeArmLambda1GB   ComputeType = 5
	ComputeTypeArmLambda2GB   ComputeType = 6
	ComputeTypeArmLambda4GB   ComputeType = 7
	ComputeTypeArmLambda8GB   ComputeType = 8
	ComputeTypeArmLambda10GB  ComputeType = 9
	ComputeTypeLinuxLambda1GB  ComputeType = 10
	ComputeTypeLinuxLambda2GB  ComputeType = 11
	ComputeTypeLinuxLambda4GB  ComputeType = 12
	ComputeTypeLinuxLambda8GB  ComputeType = 13
	ComputeTypeLinuxLambda10GB ComputeType = 14
)

type EnvironmentType uint32

const (
	EnvironmentTypeUnknown                    EnvironmentType = 0
	EnvironmentTypeLinuxContainer             EnvironmentType = 1
	EnvironmentTypeLinuxGPUContainer          EnvironmentType = 2
	EnvironmentTypeARMContainer               EnvironmentType = 3
	EnvironmentTypeWindowsServer2019Container EnvironmentType = 4
	EnvironmentTypeArmEC2Container            EnvironmentType = 5
	EnvironmentTypeLinuxEC2Container          EnvironmentType = 6
	EnvironmentTypeWindowsServer2022Container EnvironmentType = 7
	EnvironmentTypeWindowsEC2Container        EnvironmentType = 8
	EnvironmentTypeArmLambdaContainer         EnvironmentType = 9
	EnvironmentTypeLinuxLambdaContainer       EnvironmentType = 10
)
