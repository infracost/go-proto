package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FeatureGroup struct {
	resource.Resource             `tree:"-"`
	OnlineStoreEnabled            value.Bool   `tree:"online_store_enabled"`
	OnlineStoreStorageType        value.String `tree:"online_store_storage_type"`
	ThroughputMode                value.String `tree:"throughput_mode"`
	ProvisionedReadCapacityUnits  value.Int    `tree:"provisioned_read_capacity_units"`
	ProvisionedWriteCapacityUnits value.Int    `tree:"provisioned_write_capacity_units"`
}
