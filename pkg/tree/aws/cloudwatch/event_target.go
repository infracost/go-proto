package cloudwatch

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type EventTarget struct {
	resource.Resource   `tree:"-"`
	PropagateTagsToTask value.Bool   `tree:"propagate_tags_to_task"` // from task def
	TaskDefinitionID    value.String `tree:"task_definition_id"`

	Relationships EventTargetRelationships `tree:"-"`
}

type EventTargetRelationships struct {
	TaskDefinition *ecs.TaskDefinition `tree:"task_definition"`
}
