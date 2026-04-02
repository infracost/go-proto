package efs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FileSystem struct {
	resource.Resource           `tree:"-"`
	AvailabilityZoneName        value.String `tree:"availability_zone_name"`
	HasLifecyclePolicy          value.Bool   `tree:"has_lifecycle_policy"`
	ProvisionedThroughputInMiBps value.Double `tree:"provisioned_throughput_in_mibps"`
}
