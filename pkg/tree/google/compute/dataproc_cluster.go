package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DataprocCluster struct {
	resource.Resource `tree:"-"`
	WorkerMachineType value.String `tree:"worker_machine_type"`
	MasterMachineType value.String `tree:"master_machine_type"`
}
