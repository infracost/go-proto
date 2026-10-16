package eks

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Capability struct {
	resource.Resource `tree:"-"`
	Type              value.String `tree:"type"`
}
