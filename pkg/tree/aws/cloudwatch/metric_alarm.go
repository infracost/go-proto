package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MetricAlarm struct {
	resource.Resource  `tree:"-"`
	ComparisonOperator value.Value[ComparisonOperator] `tree:"comparison_operator"`
	MetricCount        value.Int                       `tree:"metric_count"`
	Period             value.Int                       `tree:"period"`
}

type ComparisonOperator uint32

const (
	ComparisonOperatorUnknown ComparisonOperator = iota
	ComparisonOperatorGreaterThanOrEqualToThreshold
	ComparisonOperatorGreaterThanThreshold
	ComparisonOperatorLessThanThreshold
	ComparisonOperatorLessThanOrEqualToThreshold
	ComparisonOperatorLessThanLowerOrGreaterThanUpperThreshold
	ComparisonOperatorLessThanLowerThreshold
	ComparisonOperatorGreaterThanUpperThreshold
)
