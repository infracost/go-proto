package cloudtrail

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Trail struct {
	resource.Resource       `tree:"-"`
	IncludeManagementEvents value.Bool `tree:"include_management_events"`
	IncludeInsightEvents    value.Bool `tree:"include_insight_events"`
}
