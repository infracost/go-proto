package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CompositeAlarm struct {
	resource.Resource `tree:"-"`
}

type ContributorInsightRule struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
}

type MetricStream struct {
	resource.Resource `tree:"-"`
}

type SyntheticsCanary struct {
	resource.Resource  `tree:"-"`
	ScheduleExpression value.String `tree:"schedule_expression"`
}

type EventArchive struct {
	resource.Resource `tree:"-"`
}

type EventAPIDestination struct {
	resource.Resource `tree:"-"`
}

type InternetMonitor struct {
	resource.Resource  `tree:"-"`
	MonitoredResources value.Int `tree:"monitored_resources"`
}

type RUMAppMonitor struct {
	resource.Resource `tree:"-"`
}
