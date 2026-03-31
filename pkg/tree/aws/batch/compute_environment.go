package batch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ComputeEnfironment struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
}
