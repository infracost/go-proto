package bedrock

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ProvisionedModelThroughput struct {
	resource.Resource  `tree:"-"`
	ModelARN           value.String `tree:"model_arn"`
	CommitmentDuration value.String `tree:"commitment_duration"`
	ModelUnits         value.Int    `tree:"model_units"`
}
