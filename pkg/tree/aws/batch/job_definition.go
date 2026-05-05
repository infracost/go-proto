package batch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

type JobDefinition struct {
	resource.Resource `tree:"-"`
}
