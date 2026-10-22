package alloydb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`     // PRIMARY | READ_POOL
	AvailabilityType  value.String `tree:"availability_type"` // ZONAL | REGIONAL
	CPUCount          value.Int    `tree:"cpu_count"`
	NodeCount         value.Int    `tree:"node_count"`
}
