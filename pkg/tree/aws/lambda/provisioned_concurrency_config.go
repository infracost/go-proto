package lambda

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ProvisionedConcurrencyConfig struct {
	resource.Resource               `tree:"-"`
	ProvisionedConcurrentExecutions value.Int `tree:"provisioned_concurrent_executions"`
}
