package glue

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DevEndpoint struct {
	resource.Resource `tree:"-"`
	NumberOfNodes     value.Int               `tree:"number_of_nodes"`
	NumberOfWorkers   value.Int               `tree:"number_of_workers"`
	WorkerType        value.Value[WorkerType] `tree:"worker_type"`
}
