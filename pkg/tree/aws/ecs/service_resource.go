package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ServiceResource struct {
	resource.Resource       `tree:"-"`
	LaunchType              value.Value[LaunchType] `tree:"launch_type"`
	DesiredCount            value.Int    `tree:"desired_count"`
	ClusterReference        value.String `tree:"cluster"`
	TaskDefinitionReference value.String `tree:"task_definition"`
	PropagateTags           value.String `tree:"propagate_tags"`
	AssignPublicIP          value.Bool   `tree:"assign_public_ip"`
	Name                    value.String `tree:"name"`

	CapacityProviderStrategy *CapacityProviderStrategy `tree:"capacity_provider_strategy"`

	Relationships ServiceRelationships `tree:"-"`
}

type LaunchType uint32

const (
	LaunchTypeEC2      LaunchType = 0
	LaunchTypeFargate  LaunchType = 1
	LaunchTypeExternal LaunchType = 2
)

type ServiceRelationships struct {
	Cluster                    *Cluster
	TaskDefinition             *TaskDefinition
	KnownCapacityProviderTypes map[string]CapacityProviderType
}
