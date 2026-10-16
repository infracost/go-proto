package waf

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type WebACLLoggingConfiguration struct {
	resource.Resource `tree:"-"`
	DestinationType   value.Value[LogDestinationType] `tree:"destination_type"`
}

type LogDestinationType uint32

const (
	LogDestinationTypeUnknown LogDestinationType = iota
	LogDestinationTypeCloudWatchLogs
	LogDestinationTypeS3
	LogDestinationTypeKinesisDataFirehose
)
