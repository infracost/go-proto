package transfer

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Server struct {
	resource.Resource `tree:"-"`
	Protocols         value.List[string] `tree:"protocols"`
}
