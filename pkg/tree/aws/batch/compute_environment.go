package batch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ComputeEnvironment struct {
	resource.Resource `tree:"-"`
	InstanceTypes     value.List[string] `tree:"instance_types"`
}
