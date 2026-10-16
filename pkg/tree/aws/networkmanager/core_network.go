package networkmanager

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CoreNetwork struct {
	resource.Resource `tree:"-"`
	EdgeLocations     value.Int `tree:"edge_locations"`
}
