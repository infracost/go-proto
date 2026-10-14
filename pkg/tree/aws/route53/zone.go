package route53

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Zone struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
}
