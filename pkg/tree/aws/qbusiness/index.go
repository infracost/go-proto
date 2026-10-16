package qbusiness

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Index struct {
	resource.Resource `tree:"-"`
	Type              value.String `tree:"type"`
	Units             value.Int    `tree:"units"`
}
