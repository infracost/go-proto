package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type WatcherFlowLog struct {
	resource.Resource                      `tree:"-"`
	Enabled                                value.Bool `tree:"enabled"`
	TrafficAnalyticsEnabled                value.Bool `tree:"traffic_analytics_enabled"`
	TrafficAnalyticsAcceleratedProcessing  value.Bool `tree:"traffic_analytics_accelerated_processing"`
}
