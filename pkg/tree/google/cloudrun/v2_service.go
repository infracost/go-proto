package cloudrun

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type V2Service struct {
	resource.Resource    `tree:"-"`
	CPULimit             value.Int    `tree:"cpu_limit"`
	MemoryLimitBytes        value.Int    `tree:"memory_limit_bytes"`
	IsThrottlingEnabled  value.Bool   `tree:"is_throttling_enabled"`
	MinInstanceCount     value.Double `tree:"min_instance_count"`
}
