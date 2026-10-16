package dynamodb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ContributorInsights struct {
	resource.Resource `tree:"-"`
	Mode              value.Value[ContributorInsightsMode] `tree:"mode"`
}

type ContributorInsightsMode uint32

const (
	ContributorInsightsModeUnknown ContributorInsightsMode = iota
	ContributorInsightsModeAccessedAndThrottledKeys
	ContributorInsightsModeThrottledKeys
)
