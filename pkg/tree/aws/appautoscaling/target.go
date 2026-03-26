package appautoscaling

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Target struct {
	resource.Resource  `tree:"-"`
	ResourceID         value.String `tree:"resource_id"`
	ScalabaleDimension value.String `tree:"scalable_dimension"`
	MinCapacity        value.Int    `tree:"min_capacity"`
	MaxCapacity        value.Int    `tree:"max_capacity"`
}
