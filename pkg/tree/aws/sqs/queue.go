package sqs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Queue struct {
	resource.Resource    `tree:"-"`
	FifoQueue            value.Bool   `tree:"fifo_queue"`
	SQSManagedSseEnabled value.Bool   `tree:"sqs_managed_sse_enabled"`
	KmsMasterKeyID       value.String `tree:"kms_master_key_id"`
}
