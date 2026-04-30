package compute

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type RegionDisk struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	SelfLink          value.String `tree:"self_link"`
	IsAttached        bool         `tree:"-"`
}
