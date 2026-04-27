package dms

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ReplicationInstance struct {
	resource.Resource        `tree:"-"`
	AllocatedStorageGB       value.Int    `tree:"allocated_storage"`
	ReplicationInstanceClass value.String `tree:"replication_instance_class"`
	MultiAZ                  value.Bool   `tree:"multi_az"`
	PubliclyAccessible       value.Bool   `tree:"publicly_accessible"`
}
