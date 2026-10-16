package apprunner

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Service struct {
	resource.Resource      `tree:"-"`
	VCPU                   value.Double `tree:"vcpu"`
	MemoryGB               value.Double `tree:"memory_gb"`
	AutoDeploymentsEnabled value.Bool   `tree:"auto_deployments_enabled"`
	FromSourceCode         value.Bool   `tree:"from_source_code"`
}
