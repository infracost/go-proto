package globalaccelerator

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EndpointGroup struct {
	resource.Resource   `tree:"-"`
	EndpointGroupRegion value.String `tree:"endpoint_group_region"`
}
