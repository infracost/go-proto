package m2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Environment struct {
	resource.Resource `tree:"-"`
	EngineType        value.Value[EngineType] `tree:"engine_type"`
	// InstanceType has the "M2." prefix removed, e.g. "m5.large".
	InstanceType    value.String `tree:"instance_type"`
	DesiredCapacity value.Int    `tree:"desired_capacity"`
}

type EngineType uint32

const (
	EngineTypeUnknown EngineType = iota
	EngineTypeMicroFocus
	EngineTypeBluAge
)
