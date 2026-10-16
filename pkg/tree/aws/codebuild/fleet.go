package codebuild

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Fleet is reserved build capacity, billed for every hour it exists.
type Fleet struct {
	resource.Resource `tree:"-"`
	BaseCapacity      value.Int    `tree:"base_capacity"`
	ComputeType       value.String `tree:"compute_type"`
	EnvironmentType   value.String `tree:"environment_type"`
	InstanceType      value.String `tree:"instance_type"`
}
