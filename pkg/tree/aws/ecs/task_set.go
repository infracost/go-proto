package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TaskSet struct {
	resource.Resource       `tree:"-"`
	ServiceReference        value.String `tree:"service"`
	TaskDefinitionReference value.String `tree:"task_definition"`

	Relationships TaskSetRelationships `tree:"-"`
}

type TaskSetRelationships struct {
	Service        *Service
	TaskDefinition *TaskDefinition
}
