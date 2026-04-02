package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MetricAlarm struct {
	resource.Resource  `tree:"-"`
	ComparisonOperator value.String `tree:"comparison_operator"`
	MetricCount        value.Int    `tree:"metric_count"`
	Period             value.Int    `tree:"period"`
}

type ComparisonOperator uint32

const (
	ComparisonOperatorGreaterThanOrEqualToThreshold            ComparisonOperator = 0
	ComparisonOperatorGreaterThanThreshold                     ComparisonOperator = 1
	ComparisonOperatorLessThanThreshold                        ComparisonOperator = 2
	ComparisonOperatorLessThanOrEqualToThreshold               ComparisonOperator = 3
	ComparisonOperatorLessThanLowerOrGreaterThanUpperThreshold ComparisonOperator = 4
	ComparisonOperatorLessThanLowerThreshold                   ComparisonOperator = 5
	ComparisonOperatorGreaterThanUpperThreshold                ComparisonOperator = 6
)
