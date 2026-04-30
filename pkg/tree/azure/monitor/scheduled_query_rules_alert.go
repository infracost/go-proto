package monitor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ScheduledQueryRulesAlert struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
	TimeSeriesCount   value.Int  `tree:"time_series_count"`
	FrequencyMinutes  value.Int  `tree:"frequency_minutes"`
}
