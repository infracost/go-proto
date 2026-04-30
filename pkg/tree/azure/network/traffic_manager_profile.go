package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TrafficManagerProfile struct {
	resource.Resource   `tree:"-"`
	Enabled             value.Bool `tree:"enabled"`
	TrafficViewEnabled  value.Bool `tree:"traffic_view_enabled"`
	HealthCheckInterval value.Int  `tree:"health_check_interval"`
}
