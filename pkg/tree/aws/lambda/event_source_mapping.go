package lambda

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EventSourceMapping struct {
	resource.Resource `tree:"-"`
	SourceType        value.Value[EventSourceType] `tree:"source_type"`
	ProvisionedMode   value.Bool                   `tree:"provisioned_mode"`
	MinimumPollers    value.Int                    `tree:"minimum_pollers"`
	PollerGroupName   value.String                 `tree:"poller_group_name"`
}

type EventSourceType uint32

const (
	EventSourceTypeUnknown EventSourceType = iota
	EventSourceTypeKafka
	EventSourceTypeSQS
)
