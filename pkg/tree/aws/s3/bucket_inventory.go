package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type BucketInventory struct {
	resource.Resource        `tree:"-"`
	Bucket                   value.String                              `tree:"bucket"`
	Name                     value.String                              `tree:"name"`
	IncludedObjectVersions   value.Value[InventoryObjectVersions]     `tree:"included_object_versions"`
	ScheduleFrequency        value.Value[InventoryScheduleFrequency]  `tree:"schedule_frequency"`
	DestinationBucketFormat  value.Value[InventoryDestinationFormat]   `tree:"destination_bucket_format"`
	DestinationBucketARN     value.String                              `tree:"destination_bucket_arn"`
}

type InventoryObjectVersions uint32

const (
	InventoryObjectVersionsUnknown InventoryObjectVersions = iota
	InventoryObjectVersionsAll
	InventoryObjectVersionsCurrent
)

type InventoryScheduleFrequency uint32

const (
	InventoryScheduleFrequencyUnknown InventoryScheduleFrequency = iota
	InventoryScheduleFrequencyDaily
	InventoryScheduleFrequencyWeekly
)

type InventoryDestinationFormat uint32

const (
	InventoryDestinationFormatUnknown InventoryDestinationFormat = iota
	InventoryDestinationFormatCSV
	InventoryDestinationFormatORC
	InventoryDestinationFormatParquet
)
