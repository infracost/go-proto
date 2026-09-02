package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// AutoscalingSchedule is a scheduled scaling action attached to an
// autoscaling group (aws_autoscaling_schedule). MinSize, MaxSize and
// DesiredCapacity use -1 to mean "do not change this value when the
// action runs", matching the AWS API.
type AutoscalingSchedule struct {
	resource.Resource    `tree:"-"`
	AutoscalingGroupName value.String `tree:"autoscaling_group_name"`
	ScheduledActionName  value.String `tree:"scheduled_action_name"`
	MinSize              value.Int    `tree:"min_size"`
	MaxSize              value.Int    `tree:"max_size"`
	DesiredCapacity      value.Int    `tree:"desired_capacity"`
	Recurrence           value.String `tree:"recurrence"`
	StartTime            value.String `tree:"start_time"`
	EndTime              value.String `tree:"end_time"`
	TimeZone             value.String `tree:"time_zone"`
}
