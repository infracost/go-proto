package kendra

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Index struct {
	resource.Resource    `tree:"-"`
	Name                 value.String `tree:"name"`
	Edition              value.String `tree:"edition"`
	QueryCapacityUnits   value.Int    `tree:"query_capacity_units"`
	StorageCapacityUnits value.Int    `tree:"storage_capacity_units"`
}
