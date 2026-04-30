package ssm

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Parameter struct {
	resource.Resource `tree:"-"`
	Tier              value.Value[ParameterTier] `tree:"tier"`
}

type ParameterTier uint32

const (
	ParameterTierUnknown            ParameterTier = iota
	ParameterTierStandard
	ParameterTierAdvanced
	ParameterTierIntelligentTiering
)
