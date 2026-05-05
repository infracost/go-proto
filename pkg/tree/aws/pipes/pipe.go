package pipes

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Pipe struct {
	resource.Resource `tree:"-"`
	TaskDefinitionID  value.String `tree:"task_definition_id"`

	Relationships PipeRelationships `tree:"-"`
}

type PipeRelationships struct {
	TaskDefinition *ecs.TaskDefinition
}
