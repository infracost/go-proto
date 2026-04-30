package glue

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Job struct {
	resource.Resource `tree:"-"`
	MaxCapacity       value.Double              `tree:"max_capacity"`
	NumberOfWorkers   value.Int                 `tree:"number_of_workers"`
	WorkerType        value.Value[WorkerType]   `tree:"worker_type"`
}

type WorkerType uint32

const (
	WorkerTypeUnknown  WorkerType = iota
	WorkerTypeStandard
	WorkerTypeG1X
	WorkerTypeG2X
	WorkerTypeG025X
	WorkerTypeG4X
	WorkerTypeG8X
	WorkerTypeZ2X
)
