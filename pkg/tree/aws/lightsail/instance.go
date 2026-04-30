package lightsail

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Instance struct {
	resource.Resource `tree:"-"`
	BundleID          value.String `tree:"bundle_id"`
}
