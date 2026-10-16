package athena

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Workgroup struct {
	resource.Resource `tree:"-"`
	EngineVersion     value.String `tree:"engine_version"`
}
