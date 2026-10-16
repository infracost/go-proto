package cloudsearch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Domain struct {
	resource.Resource `tree:"-"`
	InstanceType      value.String `tree:"instance_type"`
	PartitionCount    value.Int    `tree:"partition_count"`
	ReplicationCount  value.Int    `tree:"replication_count"`
	MultiAZ           value.Bool   `tree:"multi_az"`
}
