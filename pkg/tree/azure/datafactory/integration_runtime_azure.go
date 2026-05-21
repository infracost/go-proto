package datafactory

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type IntegrationRuntimeAzure struct {
	resource.Resource `tree:"-"`
	Cores             value.Int                   `tree:"cores"`
	ComputeType       value.Value[ADFComputeType] `tree:"compute_type"`
}
