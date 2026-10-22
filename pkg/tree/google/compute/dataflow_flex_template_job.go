package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DataflowFlexTemplateJob struct {
	resource.Resource `tree:"-"`
	MachineType       value.String `tree:"machine_type"`
	NumWorkers        value.Int    `tree:"num_workers"`
}
