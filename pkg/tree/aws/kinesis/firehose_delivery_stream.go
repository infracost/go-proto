package kinesis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type FirehoseDeliveryStream struct {
	resource.Resource           `tree:"-"`
	DataFormatConversionEnabled value.Bool   `tree:"data_format_conversion_enabled"`
	VPCDeliveryEnabled          value.Bool   `tree:"vpc_delivery_enabled"`
	VPCDeliveryAZs              value.Int    `tree:"vpc_delivery_azs"`
	ServerSideEncryption        value.Bool   `tree:"server_side_encryption"`
	KinesisStreamARN            value.String `tree:"kinesis_stream_arn"`

	Relationships FirehoseDeliveryStreamRelationships `tree:"-"`
}

type FirehoseDeliveryStreamRelationships struct {
	KinesisStream *Stream
}
