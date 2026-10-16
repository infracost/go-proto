package timestream

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type InfluxDBInstance struct {
	resource.Resource `tree:"-"`
	DBInstanceType    value.String             `tree:"db_instance_type"`
	DBStorageType     value.Value[StorageType] `tree:"db_storage_type"`
	MultiAZ           value.Bool               `tree:"multi_az"`
	AllocatedStorage  value.Int                `tree:"allocated_storage"`
}

type InfluxDBCluster struct {
	resource.Resource `tree:"-"`
	DBInstanceType    value.String             `tree:"db_instance_type"`
	DBStorageType     value.Value[StorageType] `tree:"db_storage_type"`
	AllocatedStorage  value.Int                `tree:"allocated_storage"`
	// V3 clusters use shared object storage and IO-included compute.
	V3 value.Bool `tree:"v3"`
}

type StorageType uint32

const (
	StorageTypeUnknown StorageType = iota
	StorageTypeInfluxIOIncludedT1
	StorageTypeInfluxIOIncludedT2
	StorageTypeInfluxIOIncludedT3
)
