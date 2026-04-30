package expressroute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Gateway struct {
	resource.Resource `tree:"-"`
	ScaleUnits        value.Int `tree:"scale_units"`
}
