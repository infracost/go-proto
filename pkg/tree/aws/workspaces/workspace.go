package workspaces

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workspace struct {
	resource.Resource `tree:"-"`
	ComputeType       value.Value[ComputeType] `tree:"compute_type"`
	RunningMode       value.Value[RunningMode] `tree:"running_mode"`
	RootVolumeSizeGiB value.Int                `tree:"root_volume_size_gib"`
	UserVolumeSizeGiB value.Int                `tree:"user_volume_size_gib"`
}

type ComputeType uint32

const (
	ComputeTypeUnknown ComputeType = iota
	ComputeTypeValue
	ComputeTypeStandard
	ComputeTypePerformance
	ComputeTypePower
	ComputeTypePowerPro
	ComputeTypeGraphics
	ComputeTypeGraphicsPro
	ComputeTypeGraphicsG4dn
	ComputeTypeGraphicsProG4dn
	ComputeTypeGeneralPurpose4XLarge
	ComputeTypeGeneralPurpose8XLarge
)

type RunningMode uint32

const (
	RunningModeUnknown RunningMode = iota
	RunningModeAlwaysOn
	RunningModeAutoStop
	RunningModeManual
)
