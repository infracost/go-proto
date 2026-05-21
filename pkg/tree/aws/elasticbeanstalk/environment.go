package elasticbeanstalk

import (
	"github.com/infracost/go-proto/pkg/tree/aws/cloudwatch"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/aws/rds"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Environment struct {
	resource.Resource `tree:"-"`
	Name              value.String                  `tree:"name"`
	LoadBalancerType  value.Value[LoadBalancerType] `tree:"load_balancer_type"`
	InstanceCount     value.Int                     `tree:"instance_count"`

	// Synthetic sub-resources derived from EB option_settings. EB doesn't
	// expose these as separate top-level resources, but their attributes
	// (instance type, EBS volume, DB engine, etc.) need to round-trip
	// through proto so the processor can re-hydrate the tree.
	SyntheticLaunchConfiguration *ec2.LaunchConfiguration `tree:"synthetic_launch_configuration"`
	SyntheticLoadBalancer        *ec2.LoadBalancer        `tree:"synthetic_load_balancer"`
	SyntheticClassicLoadBalancer *ec2.ClassicLoadBalancer `tree:"synthetic_classic_load_balancer"`
	SyntheticDBInstance          *rds.Instance            `tree:"synthetic_db_instance"`
	SyntheticCloudwatchLogGroup  *cloudwatch.LogGroup     `tree:"synthetic_cloudwatch_log_group"`
}

type LoadBalancerType uint32

const (
	LoadBalancerTypeUnknown LoadBalancerType = iota
	LoadBalancerTypeApplication
	LoadBalancerTypeClassic
	LoadBalancerTypeNetwork
)
