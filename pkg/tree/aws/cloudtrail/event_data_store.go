package cloudtrail

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EventDataStore struct {
	resource.Resource `tree:"-"`
	FixedRetention    value.Bool `tree:"fixed_retention"`
	RetentionDays     value.Int  `tree:"retention_days"`
}
