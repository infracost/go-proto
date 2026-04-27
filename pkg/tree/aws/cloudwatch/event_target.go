package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EventTarget struct {
	resource.Resource `tree:"-"`
	PropagateTags     value.String `tree:"propagate_tags"`
	TaskDefinitionID  value.String `tree:"task_definition_id"`
}
