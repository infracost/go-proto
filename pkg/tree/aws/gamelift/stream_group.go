package gamelift

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type StreamGroup struct {
	resource.Resource `tree:"-"`
	StreamClass       value.String          `tree:"stream_class"`
	Locations         []StreamGroupLocation `tree:"locations"`
}

type StreamGroupLocation struct {
	Name               value.String `tree:"name"`
	AlwaysOnCapacity   value.Int    `tree:"always_on_capacity"`
	TargetIdleCapacity value.Int    `tree:"target_idle_capacity"`
}
