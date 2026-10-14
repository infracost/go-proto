package cloudfront

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Function struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
}
