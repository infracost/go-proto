package dns

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Zone struct {
	resource.Resource `tree:"-"`
	ZoneType          value.Value[ZoneType] `tree:"zone_type"`
}
