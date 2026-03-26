package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Host struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
	InstanceFamily    value.String `tree:"instance_family"`
}
