package route53

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type HealthCheck struct {
	resource.Resource     `tree:"-"`
	Type                  value.Value[HealthCheckType] `tree:"type"`
	RequestIntervalSeconds value.Int                   `tree:"request_interval"`
	MeasureLatency        value.Bool                   `tree:"measure_latency"`
}

type HealthCheckType uint32

const (
	HealthCheckTypeUnknown            HealthCheckType = iota
	HealthCheckTypeHTTP
	HealthCheckTypeHTTPS
	HealthCheckTypeHTTPStrMatch
	HealthCheckTypeHTTPSStrMatch
	HealthCheckTypeTCP
	HealthCheckTypeCalculated
	HealthCheckTypeCloudWatchMetric
	HealthCheckTypeRecoveryControl
)
