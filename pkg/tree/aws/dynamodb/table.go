package dynamodb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Table struct {
	resource.Resource          `tree:"-"`
	Name                       value.String               `tree:"name"`
	BillingMode                value.Value[BillingMode]   `tree:"billing_mode"`
	WriteCapacity              value.Int                  `tree:"write_capacity"`
	ReadCapacity               value.Int                  `tree:"read_capacity"`
	PointInTimeRecoveryEnabled value.Bool                 `tree:"point_in_time_recovery_enabled"`
	ReplicaRegions             []value.String             `tree:"replica_regions"`
}

type BillingMode uint32

const (
	BillingModeProvisioned  BillingMode = 0
	BillingModePayPerRequest BillingMode = 1
)
