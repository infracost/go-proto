package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VirtualInterface struct {
	resource.Resource `tree:"-"`
	SiteLinkEnabled   value.Bool `tree:"sitelink_enabled"`
}
