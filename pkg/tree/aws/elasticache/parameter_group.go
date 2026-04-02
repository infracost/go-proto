package elasticache

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ParameterGroup struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	Family            value.String `tree:"family"`
}
