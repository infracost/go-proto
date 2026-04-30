package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LoadBalancer struct {
	resource.Resource `tree:"-"`
	Type              value.Value[LoadBalancerType] `tree:"type"`
	SubnetMappings    []SubnetMapping               `tree:"subnet_mappings"`
}

type SubnetMapping struct {
	SubnetID     value.String `tree:"subnet_id"`
	AllocationID value.String `tree:"allocation_id"`

	Relationships SubnetMappingRelationships `tree:"-"`
}

type SubnetMappingRelationships struct {
	Subnet     *Subnet
	Allocation *ElasticIP
}

type LoadBalancerType uint32

const (
	LoadBalancerTypeUnknown     LoadBalancerType = iota
	LoadBalancerTypeApplication
	LoadBalancerTypeNetwork
)
