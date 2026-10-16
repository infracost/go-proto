package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CapacityReservation struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String         `tree:"instance_type"`
	InstancePlatform  value.String         `tree:"instance_platform"`
	InstanceCount     value.Int            `tree:"instance_count"`
	Tenancy           value.Value[Tenancy] `tree:"tenancy"`
}
