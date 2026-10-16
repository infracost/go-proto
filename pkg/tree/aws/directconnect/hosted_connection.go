package directconnect

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type HostedConnection struct {
	resource.Resource `tree:"-"`
	Bandwidth         value.Value[Bandwidth] `tree:"bandwidth"`
}
