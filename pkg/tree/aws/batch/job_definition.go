package batch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type JobDefinition struct {
	resource.Resource `tree:"-"`
	PropagateTags     value.Bool `tree:"propagate_tags"`
}
