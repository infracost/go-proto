package glue

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CatalogTableOptimizer struct {
	resource.Resource `tree:"-"`
	Type              value.String `tree:"type"`
}
