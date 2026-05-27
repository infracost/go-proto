package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestTreePostProcess_IsIdempotent(t *testing.T) {
	tree := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				VPCs: []ec2.VPC{
					{Resource: resource.Resource{ID: "vpc-1"}},
				},
				VPCEndpoints: []ec2.VPCEndpoint{
					{
						Resource: resource.Resource{ID: "vpce-1"},
						VPCID:    value.New("vpc-1", 0, "", nil),
					},
				},
				Subnets: []ec2.Subnet{
					{Resource: resource.Resource{ID: "subnet-1"}},
				},
				NATGateways: []ec2.NATGateway{
					{
						Resource: resource.Resource{ID: "nat-1"},
						SubnetID: value.New("subnet-1", 0, "", nil),
					},
				},
			},
		},
	}

	tree.PostProcess()
	vpcEndpoints := append([]*ec2.VPCEndpoint(nil), tree.AWS.EC2.VPCs[0].Relationships.VPCEndpoints...)
	natGWs := append([]*ec2.NATGateway(nil), tree.AWS.EC2.Subnets[0].Relationships.NATGateways...)

	tree.PostProcess()
	assert.Equal(t, vpcEndpoints, tree.AWS.EC2.VPCs[0].Relationships.VPCEndpoints)
	assert.Equal(t, natGWs, tree.AWS.EC2.Subnets[0].Relationships.NATGateways)

	assert.Len(t, vpcEndpoints, 1)
	assert.Len(t, natGWs, 1)
}
