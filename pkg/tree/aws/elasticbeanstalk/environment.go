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

	Relationships EnvironmentRelationships `tree:"-"`
}

type EnvironmentRelationships struct {
	LoadBalancer        *ec2.LoadBalancer
	ClassicLoadBalancer *ec2.ClassicLoadBalancer
	DBInstance          *rds.Instance
	LaunchConfiguration *ec2.LaunchConfiguration
	CloudwatchLogGroup  *cloudwatch.LogGroup
}

type LoadBalancerType uint32

const (
	LoadBalancerTypeUnknown LoadBalancerType = iota
	LoadBalancerTypeApplication
	LoadBalancerTypeClassic
	LoadBalancerTypeNetwork
)
