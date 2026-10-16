package sesv2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// DedicatedIPAssignment is one leased standard dedicated IP.
type DedicatedIPAssignment struct {
	resource.Resource `tree:"-"`
}

type DedicatedIPPool struct {
	resource.Resource `tree:"-"`
	ScalingMode       value.Value[ScalingMode] `tree:"scaling_mode"`
}

type ScalingMode uint32

const (
	ScalingModeUnknown ScalingMode = iota
	ScalingModeStandard
	ScalingModeManaged
)

type EmailIdentity struct {
	resource.Resource `tree:"-"`
}

type ReceiptRule struct {
	resource.Resource `tree:"-"`
}
