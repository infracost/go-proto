package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FlowLog struct {
	resource.Resource  `tree:"-"`
	LogDestinationType value.Value[FlowLogDestinationType] `tree:"log_destination_type"`
	ParquetFormat      value.Bool                          `tree:"parquet_format"`
}

type FlowLogDestinationType uint32

const (
	FlowLogDestinationTypeUnknown FlowLogDestinationType = iota
	FlowLogDestinationTypeCloudWatchLogs
	FlowLogDestinationTypeS3
	FlowLogDestinationTypeKinesisDataFirehose
)
