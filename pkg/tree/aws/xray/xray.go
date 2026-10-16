package xray

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type XRay struct {
	Groups        []Group        `tree:"groups"`
	SamplingRules []SamplingRule `tree:"sampling_rules"`
}

type Group struct {
	resource.Resource `tree:"-"`
	InsightsEnabled   value.Bool `tree:"insights_enabled"`
}

type SamplingRule struct {
	resource.Resource `tree:"-"`
}
