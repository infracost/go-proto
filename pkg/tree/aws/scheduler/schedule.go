package scheduler

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Schedule struct {
	resource.Resource `tree:"-"`
	TaskDefinitionARN value.String `tree:"task_definition_arn"`

	Relationships ScheduleRelationships `tree:"-"`
}

type ScheduleRelationships struct {
	TaskDefinition *ecs.TaskDefinition
}
