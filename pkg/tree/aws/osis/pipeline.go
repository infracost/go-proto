package osis

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Pipeline struct {
	resource.Resource `tree:"-"`
	MinUnits          value.Int `tree:"min_units"`
}
