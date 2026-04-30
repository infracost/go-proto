package kinesis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Stream struct {
	resource.Resource `tree:"-"`
	StreamMode        value.Value[StreamMode]         `tree:"stream_mode"`
	ShardCount        value.Int                       `tree:"shard_count"`
	EncryptionType    value.Value[StreamEncryptionType] `tree:"encryption_type"`
	KmsKeyID          value.String                    `tree:"kms_key_id"`
}

type StreamMode uint32

const (
	StreamModeUnknown     StreamMode = iota
	StreamModeOnDemand
	StreamModeProvisioned
)

type StreamEncryptionType uint32

const (
	StreamEncryptionTypeUnknown StreamEncryptionType = iota
	StreamEncryptionTypeNone
	StreamEncryptionTypeKMS
)
