package glue

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MLTransform struct {
	resource.Resource `tree:"-"`
	MaxCapacity       value.Double            `tree:"max_capacity"`
	NumberOfWorkers   value.Int               `tree:"number_of_workers"`
	WorkerType        value.Value[WorkerType] `tree:"worker_type"`
}
