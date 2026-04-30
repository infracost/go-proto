package cloudrun

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type V2Job struct {
	resource.Resource `tree:"-"`
	CPULimit          value.Int `tree:"cpu_limit"`
	MemoryLimitBytes     value.Int `tree:"memory_limit_bytes"`
	TaskCount         value.Int `tree:"task_count"`
}
