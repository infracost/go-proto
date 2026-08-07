package tree

import (
	"sort"
	"testing"

	"github.com/infracost/go-proto/pkg/address"
	"github.com/infracost/go-proto/pkg/tree/aws"
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/aws/dynamodb"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// res builds a resource.Resource with an ID and a parseable address so it can
// be indexed and matched by Components.
func res(id, addr string) resource.Resource {
	return resource.Resource{
		ID:         id,
		Definition: resource.Definition{Address: address.ParseOrLiteral(addr)},
	}
}

// componentAddressSets returns each component's addresses as a sorted slice, so
// assertions don't depend on within-component ordering.
func componentAddressSets(comps []Component) [][]string {
	out := make([][]string, len(comps))
	for i, c := range comps {
		addrs := c.Addresses()
		sort.Strings(addrs)
		out[i] = addrs
	}
	return out
}

func TestComponents_Empty(t *testing.T) {
	tr := &Tree{}
	assert.Empty(t, tr.Components())
}

func TestComponents_IsolatedResourcesAreSingletons(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: res("i-1", "aws_instance.a")},
					{Resource: res("i-2", "aws_instance.b")},
					{Resource: res("i-3", "aws_instance.c")},
				},
			},
		},
	}
	tr.PostProcess()

	comps := tr.Components()
	require.Len(t, comps, 3)
	// Order follows ToResources order, and each is a singleton.
	assert.Equal(t, [][]string{
		{"aws_instance.a"},
		{"aws_instance.b"},
		{"aws_instance.c"},
	}, componentAddressSets(comps))
}

// Pointer relationship: Instance.Relationships.InstanceState (*T), linked by
// PostProcess when InstanceID matches the instance ID.
func TestComponents_PointerRelationship(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: res("i-1", "aws_instance.a")},
					{Resource: res("i-2", "aws_instance.b")},
				},
				InstanceStates: []ec2.InstanceStateMapping{
					{
						Resource:   res("is-1", "aws_ec2_instance_state.a"),
						InstanceID: value.New("i-1", 0, "", nil),
					},
				},
			},
		},
	}
	tr.PostProcess()
	require.NotNil(t, tr.AWS.EC2.Instances[0].Relationships.InstanceState, "precondition: state linked")

	comps := tr.Components()
	require.Len(t, comps, 2)
	assert.ElementsMatch(t, [][]string{
		{"aws_ec2_instance_state.a", "aws_instance.a"}, // coupled
		{"aws_instance.b"},                             // independent
	}, componentAddressSets(comps))
}

// Slice-of-pointer relationship: VPC.Relationships.VPCEndpoints ([]*T), linked
// by PostProcess when the endpoint's VPCID matches the VPC ID.
func TestComponents_SlicePointerRelationship(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				VPCs: []ec2.VPC{
					{Resource: res("vpc-1", "aws_vpc.a")},
				},
				VPCEndpoints: []ec2.VPCEndpoint{
					{Resource: res("vpce-1", "aws_vpc_endpoint.a"), VPCID: value.New("vpc-1", 0, "", nil)},
					{Resource: res("vpce-2", "aws_vpc_endpoint.b"), VPCID: value.New("vpc-1", 0, "", nil)},
				},
			},
		},
	}
	tr.PostProcess()
	require.Len(t, tr.AWS.EC2.VPCs[0].Relationships.VPCEndpoints, 2, "precondition: endpoints linked")

	comps := tr.Components()
	require.Len(t, comps, 1)
	assert.Equal(t, []string{"aws_vpc.a", "aws_vpc_endpoint.a", "aws_vpc_endpoint.b"}, componentAddressSets(comps)[0])
}

// Value-slice relationship: DynamoDB Table.Relationships.AppAutoscalingTargets
// ([]T — stored as value copies, not pointers). This is the case that pointer
// identity would miss: the linked target is a *copy* of the canonical resource,
// so Components must match on address. PostProcess links by matching the
// target's ResourceID to "table/<name>".
func TestComponents_ValueSliceRelationship_MatchesByAddress(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			DynamoDB: dynamodb.DynamoDB{
				Tables: []dynamodb.Table{
					{Resource: res("tbl-1", "aws_dynamodb_table.a"), Name: value.New("mytable", 0, "name", nil)},
				},
			},
			AppAutoScaling: appautoscaling.AppAutoScaling{
				Targets: []appautoscaling.Target{
					{
						Resource:   res("tgt-1", "aws_appautoscaling_target.a"),
						ResourceID: value.New("table/mytable", 0, "resource_id", nil),
					},
				},
			},
		},
	}
	tr.PostProcess()
	require.Len(t, tr.AWS.DynamoDB.Tables[0].Relationships.AppAutoscalingTargets, 1, "precondition: target linked")
	// Sanity: the linked target is a value copy, not the canonical pointer.
	assert.NotSame(t, &tr.AWS.AppAutoScaling.Targets[0], &tr.AWS.DynamoDB.Tables[0].Relationships.AppAutoscalingTargets[0])

	comps := tr.Components()
	require.Len(t, comps, 1)
	assert.Equal(t, []string{"aws_appautoscaling_target.a", "aws_dynamodb_table.a"}, componentAddressSets(comps)[0])
}

// Transitivity: A links B and B links C should all collapse into one component,
// even though A and C are never directly linked. Built with a VPC→endpoint and
// an endpoint→subnet link sharing the middle endpoint.
func TestComponents_Transitive(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				VPCs: []ec2.VPC{
					{Resource: res("vpc-1", "aws_vpc.a")},
				},
				Subnets: []ec2.Subnet{
					{Resource: res("subnet-1", "aws_subnet.a")},
				},
				VPCEndpoints: []ec2.VPCEndpoint{
					{
						Resource:  res("vpce-1", "aws_vpc_endpoint.a"),
						VPCID:     value.New("vpc-1", 0, "", nil),
						SubnetIDs: *value.NewList([]value.String{value.New("subnet-1", 0, "", nil)}, 0, "", nil),
					},
				},
			},
		},
	}
	tr.PostProcess()
	// Precondition: endpoint linked to both vpc and subnet.
	require.Len(t, tr.AWS.EC2.VPCs[0].Relationships.VPCEndpoints, 1)
	require.Len(t, tr.AWS.EC2.VPCEndpoints[0].Relationships.Subnets, 1)

	comps := tr.Components()
	require.Len(t, comps, 1, "vpc—endpoint—subnet should collapse into one component")
	assert.Equal(t, []string{"aws_subnet.a", "aws_vpc.a", "aws_vpc_endpoint.a"}, componentAddressSets(comps)[0])
}

// A resource with no address can't be matched, so it forms its own singleton
// and must not panic.
func TestComponents_NilAddressIsSingleton(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: resource.Resource{ID: "i-1"}}, // no Definition.Address
					{Resource: res("i-2", "aws_instance.b")},
				},
			},
		},
	}
	tr.PostProcess()

	comps := tr.Components()
	require.Len(t, comps, 2)
	// The addressless resource yields an empty address set; the other is named.
	sets := componentAddressSets(comps)
	assert.ElementsMatch(t, [][]string{{}, {"aws_instance.b"}}, sets)
}

// Every resource must appear in exactly one component (a partition), and the
// component count times membership must equal the resource count.
func TestComponents_PartitionsAllResources(t *testing.T) {
	tr := &Tree{
		AWS: aws.AWS{
			EC2: ec2.EC2{
				Instances: []ec2.Instance{
					{Resource: res("i-1", "aws_instance.a")},
					{Resource: res("i-2", "aws_instance.b")},
				},
				InstanceStates: []ec2.InstanceStateMapping{
					{Resource: res("is-1", "aws_ec2_instance_state.a"), InstanceID: value.New("i-1", 0, "", nil)},
				},
				VPCs:         []ec2.VPC{{Resource: res("vpc-1", "aws_vpc.a")}},
				VPCEndpoints: []ec2.VPCEndpoint{{Resource: res("vpce-1", "aws_vpc_endpoint.a"), VPCID: value.New("vpc-1", 0, "", nil)}},
			},
		},
	}
	tr.PostProcess()

	all := tr.ToResources(false)
	comps := tr.Components()

	seen := map[string]int{}
	total := 0
	for _, c := range comps {
		total += len(c)
		for _, r := range c {
			seen[r.GetBase().ID]++
		}
	}
	assert.Equal(t, len(all), total, "components should cover every resource exactly once")
	for id, count := range seen {
		assert.Equalf(t, 1, count, "resource %s appeared in %d components", id, count)
	}
}

func TestComponent_Addresses_SkipsNilAddress(t *testing.T) {
	c := Component{
		&ec2.Instance{Resource: res("i-1", "aws_instance.a")},
		&ec2.Instance{Resource: resource.Resource{ID: "i-2"}}, // nil address
	}
	assert.Equal(t, []string{"aws_instance.a"}, c.Addresses())
}
