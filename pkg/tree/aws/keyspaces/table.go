package keyspaces

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Table struct {
	resource.Resource          `tree:"-"`
	ThroughputMode             value.Value[ThroughputMode] `tree:"throughput_mode"`
	ReadCapacityUnits          value.Int                   `tree:"read_capacity_units"`
	WriteCapacityUnits         value.Int                   `tree:"write_capacity_units"`
	PointInTimeRecoveryEnabled value.Bool                  `tree:"point_in_time_recovery_enabled"`
	TTLEnabled                 value.Bool                  `tree:"ttl_enabled"`
}

type ThroughputMode uint32

const (
	ThroughputModeUnknown ThroughputMode = iota
	ThroughputModePayPerRequest
	ThroughputModeProvisioned
)
