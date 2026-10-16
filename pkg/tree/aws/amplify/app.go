package amplify

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type App struct {
	resource.Resource `tree:"-"`
	BuildComputeType  value.Value[BuildComputeType] `tree:"build_compute_type"`
}

type BuildComputeType uint32

const (
	BuildComputeTypeUnknown BuildComputeType = iota
	BuildComputeTypeStandard8GB
	BuildComputeTypeLarge16GB
	BuildComputeTypeXLarge72GB
)
