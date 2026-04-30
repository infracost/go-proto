package network

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LoadBalancerRule struct {
	resource.Resource `tree:"-"`
	LoadBalancerID    value.String `tree:"load_balancer_id"`

	Relationships LoadBalancerRuleRelationships `tree:"-"`
}

type LoadBalancerRuleRelationships struct {
	LoadBalancer *LoadBalancer
}
