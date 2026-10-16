package msk

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Connector struct {
	resource.Resource `tree:"-"`
	MCUCount          value.Int  `tree:"mcu_count"`
	WorkerCount       value.Int  `tree:"worker_count"`
	Autoscaling       value.Bool `tree:"autoscaling"`
}
