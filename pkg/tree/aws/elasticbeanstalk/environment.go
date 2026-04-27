package elasticbeanstalk

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Environment struct {
	resource.Resource `tree:"-"`
	Name              value.String                   `tree:"name"`
	LoadBalancerType  value.Value[LoadBalancerType]   `tree:"load_balancer_type"`
	InstanceCount     value.Int                      `tree:"instance_count"`
}

type LoadBalancerType uint32

const (
	LoadBalancerTypeApplication LoadBalancerType = 0
	LoadBalancerTypeClassic     LoadBalancerType = 1
	LoadBalancerTypeNetwork     LoadBalancerType = 2
)
