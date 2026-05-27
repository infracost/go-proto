package ec2

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	ec2 := &EC2{
		Instances: []Instance{
			{
				Resource:         resource.Resource{ID: "i-1"},
				LaunchTemplateID: value.New("lt-1", 0, "", nil),
			},
		},
		InstanceStates: []InstanceStateMapping{
			{InstanceID: value.New("i-1", 0, "", nil)},
		},
		LaunchTemplates: []LaunchTemplate{
			{Resource: resource.Resource{ID: "lt-1"}, Name: value.New("lt-name", 0, "", nil)},
		},
		AutoscalingGroups: []AutoscalingGroup{
			{
				Resource:                      resource.Resource{ID: "asg-1"},
				LaunchTemplateID:              value.New("lt-1", 0, "", nil),
				MixedInstanceLaunchTemplateID: value.New("lt-1", 0, "", nil),
				LaunchConfigurationName:       value.New("lc-1", 0, "", nil),
			},
		},
		LaunchConfigurations: []LaunchConfiguration{
			{Resource: resource.Resource{ID: "lc-1"}, Name: value.New("lc-1", 0, "", nil)},
		},
		VPCs: []VPC{
			{Resource: resource.Resource{ID: "vpc-1"}},
		},
		Subnets: []Subnet{
			{Resource: resource.Resource{ID: "subnet-1"}},
		},
		VPCEndpoints: []VPCEndpoint{
			{
				Resource: resource.Resource{ID: "vpce-1"},
				VPCID:    value.New("vpc-1", 0, "", nil),
				SubnetIDs: *value.NewList([]value.Value[string]{
					value.New("subnet-1", 0, "", nil),
				}, 0, "", nil),
			},
		},
		NATGateways: []NATGateway{
			{
				Resource:     resource.Resource{ID: "nat-1"},
				SubnetID:     value.New("subnet-1", 0, "", nil),
				AllocationID: value.New("eipalloc-1", 0, "", nil),
			},
		},
		ElasticIPs: []ElasticIP{
			{Resource: resource.Resource{ID: "eipalloc-1"}},
		},
		ElasticIPAssociations: []ElasticIPAssociation{
			{
				Resource:     resource.Resource{ID: "eipassoc-1"},
				AllocationID: value.New("eipalloc-1", 0, "", nil),
			},
		},
		TransitGateways: []TransitGateway{
			{Resource: resource.Resource{ID: "tgw-1"}},
		},
		TransitGatewayVPCAttachments: []TransitGatewayVPCAttachment{
			{
				Resource:         resource.Resource{ID: "tgw-att-1"},
				VPCID:            value.New("vpc-1", 0, "", nil),
				TransitGatewayID: value.New("tgw-1", 0, "", nil),
			},
		},
		TransitGatewayPeeringAttachments: []TransitGatewayPeeringAttachment{
			{
				Resource:         resource.Resource{ID: "tgw-peer-1"},
				TransitGatewayID: value.New("tgw-1", 0, "", nil),
			},
		},
		LoadBalancers: []LoadBalancer{
			{
				Resource: resource.Resource{ID: "lb-1"},
				SubnetMappings: []SubnetMapping{
					{
						AllocationID: value.New("eipalloc-1", 0, "", nil),
						SubnetID:     value.New("subnet-1", 0, "", nil),
					},
				},
			},
		},
	}

	ec2.PostProcess()

	snapshot := struct {
		instanceState      *InstanceStateMapping
		instanceLT         *LaunchTemplate
		asgLT              *LaunchTemplate
		asgMixedLT         *LaunchTemplate
		asgLC              *LaunchConfiguration
		vpcEndpoints       []*VPCEndpoint
		vpcEndpointVPC     *VPC
		vpcEndpointSubnets []*Subnet
		subnetNATGWs       []*NATGateway
		natGWSubnet        *Subnet
		natGWEIP           *ElasticIP
		eipAssociation     *ElasticIPAssociation
		eipNATGW           *NATGateway
		eipLB              *LoadBalancer
		associationEIP     *ElasticIP
		tgwAttVPC          *VPC
		tgwAttTGW          *TransitGateway
		tgwPeerTGW         *TransitGateway
		lbMappingAlloc     *ElasticIP
		lbMappingSubnet    *Subnet
	}{
		instanceState:      ec2.Instances[0].Relationships.InstanceState,
		instanceLT:         ec2.Instances[0].Relationships.LaunchTemplate,
		asgLT:              ec2.AutoscalingGroups[0].Relationships.LaunchTemplate,
		asgMixedLT:         ec2.AutoscalingGroups[0].Relationships.MixedInstanceLaunchTemplate,
		asgLC:              ec2.AutoscalingGroups[0].Relationships.LaunchConfiguration,
		vpcEndpoints:       append([]*VPCEndpoint(nil), ec2.VPCs[0].Relationships.VPCEndpoints...),
		vpcEndpointVPC:     ec2.VPCEndpoints[0].Relationships.VPC,
		vpcEndpointSubnets: append([]*Subnet(nil), ec2.VPCEndpoints[0].Relationships.Subnets...),
		subnetNATGWs:       append([]*NATGateway(nil), ec2.Subnets[0].Relationships.NATGateways...),
		natGWSubnet:        ec2.NATGateways[0].Relationships.Subnet,
		natGWEIP:           ec2.NATGateways[0].Relationships.AllocatedElasticIP,
		eipAssociation:     ec2.ElasticIPs[0].Relationships.Association,
		eipNATGW:           ec2.ElasticIPs[0].Relationships.NATGateway,
		eipLB:              ec2.ElasticIPs[0].Relationships.LoadBalancer,
		associationEIP:     ec2.ElasticIPAssociations[0].Relationships.ElasticIP,
		tgwAttVPC:          ec2.TransitGatewayVPCAttachments[0].Relationships.VPC,
		tgwAttTGW:          ec2.TransitGatewayVPCAttachments[0].Relationships.TransitGateway,
		tgwPeerTGW:         ec2.TransitGatewayPeeringAttachments[0].Relationships.TransitGateway,
		lbMappingAlloc:     ec2.LoadBalancers[0].SubnetMappings[0].Relationships.Allocation,
		lbMappingSubnet:    ec2.LoadBalancers[0].SubnetMappings[0].Relationships.Subnet,
	}

	ec2.PostProcess()

	assert.Equal(t, snapshot.instanceState, ec2.Instances[0].Relationships.InstanceState)
	assert.Equal(t, snapshot.instanceLT, ec2.Instances[0].Relationships.LaunchTemplate)
	assert.Equal(t, snapshot.asgLT, ec2.AutoscalingGroups[0].Relationships.LaunchTemplate)
	assert.Equal(t, snapshot.asgMixedLT, ec2.AutoscalingGroups[0].Relationships.MixedInstanceLaunchTemplate)
	assert.Equal(t, snapshot.asgLC, ec2.AutoscalingGroups[0].Relationships.LaunchConfiguration)
	assert.Equal(t, snapshot.vpcEndpoints, ec2.VPCs[0].Relationships.VPCEndpoints)
	assert.Equal(t, snapshot.vpcEndpointVPC, ec2.VPCEndpoints[0].Relationships.VPC)
	assert.Equal(t, snapshot.vpcEndpointSubnets, ec2.VPCEndpoints[0].Relationships.Subnets)
	assert.Equal(t, snapshot.subnetNATGWs, ec2.Subnets[0].Relationships.NATGateways)
	assert.Equal(t, snapshot.natGWSubnet, ec2.NATGateways[0].Relationships.Subnet)
	assert.Equal(t, snapshot.natGWEIP, ec2.NATGateways[0].Relationships.AllocatedElasticIP)
	assert.Equal(t, snapshot.eipAssociation, ec2.ElasticIPs[0].Relationships.Association)
	assert.Equal(t, snapshot.eipNATGW, ec2.ElasticIPs[0].Relationships.NATGateway)
	assert.Equal(t, snapshot.eipLB, ec2.ElasticIPs[0].Relationships.LoadBalancer)
	assert.Equal(t, snapshot.associationEIP, ec2.ElasticIPAssociations[0].Relationships.ElasticIP)
	assert.Equal(t, snapshot.tgwAttVPC, ec2.TransitGatewayVPCAttachments[0].Relationships.VPC)
	assert.Equal(t, snapshot.tgwAttTGW, ec2.TransitGatewayVPCAttachments[0].Relationships.TransitGateway)
	assert.Equal(t, snapshot.tgwPeerTGW, ec2.TransitGatewayPeeringAttachments[0].Relationships.TransitGateway)
	assert.Equal(t, snapshot.lbMappingAlloc, ec2.LoadBalancers[0].SubnetMappings[0].Relationships.Allocation)
	assert.Equal(t, snapshot.lbMappingSubnet, ec2.LoadBalancers[0].SubnetMappings[0].Relationships.Subnet)
}
