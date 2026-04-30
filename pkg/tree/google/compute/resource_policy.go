package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ResourcePolicy struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	MaxRetentionDays  value.Double `tree:"max_retention_days"`
	HasSchedule       value.Bool   `tree:"has_schedule"`
}
