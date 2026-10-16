package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Fleet is a spot fleet request or an EC2 fleet. Capacities are in target
// capacity units. InstanceType and WeightedCapacity come from the first
// launch specification or override that sets an instance type, else from
// the linked launch template.
type Fleet struct {
	resource.Resource  `tree:"-"`
	InstanceType       value.String `tree:"instance_type"`
	WeightedCapacity   value.Double `tree:"weighted_capacity"`
	OnDemandCapacity   value.Int    `tree:"on_demand_capacity"`
	SpotCapacity       value.Int    `tree:"spot_capacity"`
	LaunchTemplateID   value.String `tree:"launch_template_id"`
	LaunchTemplateName value.String `tree:"launch_template_name"`

	Relationships FleetRelationships `tree:"-"`
}

type FleetRelationships struct {
	LaunchTemplate *LaunchTemplate
}

// linkLaunchTemplate links the fleet's launch template and takes its
// instance type when no override sets one.
func (f *Fleet) linkLaunchTemplate(templates []LaunchTemplate) {
	f.Relationships = FleetRelationships{}
	if f.LaunchTemplateID.IsEmpty() && f.LaunchTemplateName.IsEmpty() {
		return
	}
	for i := range templates {
		lt := &templates[i]
		if (!f.LaunchTemplateID.IsEmpty() && (f.LaunchTemplateID.Equal(lt.ID) || f.LaunchTemplateID.Value() == lt.Name.Value())) ||
			(!f.LaunchTemplateName.IsEmpty() && f.LaunchTemplateName.Value() == lt.Name.Value()) {
			f.Relationships.LaunchTemplate = lt
			return
		}
	}
}

// EffectiveInstanceType is the fleet's own instance type, else the launch template's.
func (f *Fleet) EffectiveInstanceType() string {
	if !f.InstanceType.IsEmpty() {
		return f.InstanceType.Value()
	}
	if lt := f.Relationships.LaunchTemplate; lt != nil {
		return lt.InstanceType.Value()
	}
	return ""
}
