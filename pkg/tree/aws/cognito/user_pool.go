package cognito

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type UserPool struct {
	resource.Resource `tree:"-"`
	Tier              value.Value[Tier] `tree:"tier"`
}

type Tier uint32

const (
	TierUnknown Tier = iota
	TierLite
	TierEssentials
	TierPlus
)
