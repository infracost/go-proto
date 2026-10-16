package bedrock

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CustomModel struct {
	resource.Resource   `tree:"-"`
	BaseModelIdentifier value.String `tree:"base_model_identifier"`
}
