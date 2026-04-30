package monitor

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MetricAlert struct {
	resource.Resource              `tree:"-"`
	Enabled                        value.Bool `tree:"enabled"`
	ScopeCount                     value.Int  `tree:"scope_count"`
	CriteriaDimensionsCount        value.Int  `tree:"criteria_dimensions_count"`
	DynamicCriteriaDimensionsCount value.Int  `tree:"dynamic_criteria_dimensions_count"`
}
