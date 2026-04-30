package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type AutoscalingGroup struct {
	resource.Resource                 `tree:"-"`
	Name                              value.String        `tree:"name"`
	LaunchTemplateID                  value.String        `tree:"launch_template_id"`
	LaunchTemplateName                value.String        `tree:"launch_template_name"`
	LaunchConfigurationName           value.String        `tree:"launch_configuration_name"`
	MixedInstanceLaunchTemplateID     value.String        `tree:"mixed_instance_launch_template_id"`
	MixedInstancePolicy               MixedInstancePolicy `tree:"mixed_instance_policy"`
	DesiredCapacity                   value.Int           `tree:"desired_capacity"`
	MinSize                           value.Int           `tree:"min_size"`

	Relationships AutoscalingGroupRelationships `tree:"-"`
}

type AutoscalingGroupRelationships struct {
	LaunchTemplate              *LaunchTemplate
	MixedInstanceLaunchTemplate *LaunchTemplate
	LaunchConfiguration         *LaunchConfiguration
}

type MixedInstancePolicy struct {
	OnDemandBaseCapacity                   value.Int    `tree:"on_demand_base_capacity"`
	OnDemandPercentageAboveBaseCapacity    value.Int    `tree:"on_demand_percentage_above_base_capacity"`
	LaunchTemplateOverrideInstanceType     value.String `tree:"launch_template_override_instance_type"`
	LaunchTemplateOverrideWeightedCapacity value.Int    `tree:"launch_template_override_weighted_capacity"`
}
