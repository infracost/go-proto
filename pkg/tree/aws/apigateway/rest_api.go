package apigateway

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type RestAPI struct {
	resource.Resource `tree:"-"`
	ResourceID        value.String `tree:"resource_id"`
	ScalableDimension value.String `tree:"scalable_dimension"`
	MinCapacity       value.Int    `tree:"min_capacity"`
	MaxCamacity       value.Int    `tree:"max_capacity"`
}
