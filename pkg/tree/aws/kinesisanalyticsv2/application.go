package kinesisanalyticsv2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Application struct {
	resource.Resource  `tree:"-"`
	RuntimeEnvironment value.Value[RuntimeEnvironment] `tree:"runtime_environment"`
}

type RuntimeEnvironment uint32

const (
	RuntimeEnvironmentUnknown   RuntimeEnvironment = iota
	RuntimeEnvironmentSQL1_0
	RuntimeEnvironmentFlink1_6
	RuntimeEnvironmentFlink1_8
	RuntimeEnvironmentFlink1_11
	RuntimeEnvironmentFlink1_13
	RuntimeEnvironmentFlink1_15
	RuntimeEnvironmentFlink1_18
	RuntimeEnvironmentFlink1_19
	RuntimeEnvironmentFlink1_20
)
