package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DefaultSecurityGroup struct {
	resource.Resource `tree:"-"`
	VPCID             value.String        `tree:"vpc_id"`
	IngressRules      []SecurityGroupRule `tree:"ingress"`
	EgressRules       []SecurityGroupRule `tree:"egress"`
}

type SecurityGroupRule struct {
	FromPort   value.Int                          `tree:"from_port"`
	ToPort     value.Int                          `tree:"to_port"`
	Protocol   value.Value[SecurityGroupProtocol] `tree:"protocol"`
	CIDRBlocks value.List[string]                 `tree:"cidr_blocks"`
}

type SecurityGroupProtocol uint32

const (
	SecurityGroupProtocolUnknown SecurityGroupProtocol = iota
	SecurityGroupProtocolTCP
	SecurityGroupProtocolUDP
	SecurityGroupProtocolICMP
	SecurityGroupProtocolICMPV6
	SecurityGroupProtocolAll
)
