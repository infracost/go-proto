package loganalytics

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workspace struct {
	resource.Resource            `tree:"-"`
	SKU                          value.String `tree:"sku"`
	ReservationCapacityInGBPerDay value.Int    `tree:"reservation_capacity_in_gb_per_day"`
	RetentionInDays              value.Int    `tree:"retention_in_days"`
	SentinelEnabled              value.Bool   `tree:"sentinel_enabled"`
}
