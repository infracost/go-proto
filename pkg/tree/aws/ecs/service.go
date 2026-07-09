package ecs

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Service struct {
	resource.Resource       `tree:"-"`
	LaunchType              value.Value[LaunchType] `tree:"launch_type"`
	PlatformVersion         value.String            `tree:"platform_version"`
	DesiredCount            value.Int               `tree:"desired_count"`
	ClusterReference        value.String            `tree:"cluster"`
	TaskDefinitionReference value.String            `tree:"task_definition"`
	AssignPublicIP          value.Bool              `tree:"assign_public_ip"`
	Name                    value.String            `tree:"name"`
	SubnetIDs               value.List[string]      `tree:"subnet_ids"`

	CapacityProviderStrategies []CapacityProviderStrategy `tree:"capacity_provider_strategy"`

	Relationships ServiceRelationships `tree:"-"`
}

type LaunchType uint32

const (
	LaunchTypeUnknown LaunchType = iota
	LaunchTypeEC2
	LaunchTypeFargate
	LaunchTypeExternal
)

type ServiceRelationships struct {
	Cluster                    *Cluster
	TaskDefinition             *TaskDefinition
	KnownCapacityProviderTypes map[string]CapacityProviderType
	Subnets                    []*ec2.Subnet
}
