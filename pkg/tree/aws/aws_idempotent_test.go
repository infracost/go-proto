package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudwatch"
	"github.com/infracost/go-proto/pkg/tree/aws/directconnect"
	"github.com/infracost/go-proto/pkg/tree/aws/dynamodb"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/aws/eks"
	"github.com/infracost/go-proto/pkg/tree/aws/elasticache"
	"github.com/infracost/go-proto/pkg/tree/aws/msk"
	"github.com/infracost/go-proto/pkg/tree/aws/pipes"
	"github.com/infracost/go-proto/pkg/tree/aws/scheduler"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestAWSPostProcess_IsIdempotent(t *testing.T) {
	a := &AWS{
		AppAutoScaling: appautoscaling.AppAutoScaling{
			Targets: []appautoscaling.Target{
				{
					Resource:   resource.Resource{ID: "tgt-msk"},
					ResourceID: value.New("msk-cluster-1", 0, "", nil),
				},
				{
					Resource:   resource.Resource{ID: "tgt-ddb"},
					ResourceID: value.New("table/my-table", 0, "", nil),
				},
				{
					Resource:   resource.Resource{ID: "tgt-rg"},
					ResourceID: value.New("replication-group/rg-1", 0, "", nil),
				},
			},
		},
		MSK: msk.MSK{
			Clusters: []msk.Cluster{
				{Resource: resource.Resource{ID: "msk-cluster-1"}},
			},
		},
		DynamoDB: dynamodb.DynamoDB{
			Tables: []dynamodb.Table{
				{
					Resource: resource.Resource{ID: "ddb-1"},
					Name:     value.New("my-table", 0, "", nil),
				},
			},
		},
		ElastiCache: elasticache.ElastiCache{
			ReplicationGroups: []elasticache.ReplicationGroup{
				{
					Resource: resource.Resource{ID: "rg-1"},
					ID:       value.New("rg-1", 0, "", nil),
				},
			},
		},
		EKS: eks.EKS{
			NodeGroups: []eks.NodeGroup{
				{
					Resource:         resource.Resource{ID: "ng-1"},
					LaunchTemplateID: value.New("lt-1", 0, "", nil),
				},
			},
		},
		EC2: ec2.EC2{
			LaunchTemplates: []ec2.LaunchTemplate{
				{Resource: resource.Resource{ID: "lt-1"}},
			},
			Subnets: []ec2.Subnet{
				{Resource: resource.Resource{ID: "subnet-1"}},
			},
			TransitGateways: []ec2.TransitGateway{
				{Resource: resource.Resource{ID: "tgw-1"}},
			},
		},
		Scheduler: scheduler.Scheduler{
			Schedules: []scheduler.Schedule{
				{
					Resource:          resource.Resource{ID: "sch-1"},
					TaskDefinitionARN: value.New("td-1", 0, "", nil),
				},
			},
		},
		CloudWatch: cloudwatch.CloudWatch{
			EventTargets: []cloudwatch.EventTarget{
				{
					Resource:         resource.Resource{ID: "et-1"},
					TaskDefinitionID: value.New("td-1", 0, "", nil),
				},
			},
		},
		Pipes: pipes.Pipes{
			Pipes: []pipes.Pipe{
				{
					Resource:         resource.Resource{ID: "pipe-1"},
					TaskDefinitionID: value.New("td-1", 0, "", nil),
				},
			},
		},
		ECS: ecs.ECS{
			TaskDefinitions: []ecs.TaskDefinition{
				{Resource: resource.Resource{ID: "td-1"}},
			},
			Services: []ecs.Service{
				{
					Resource: resource.Resource{ID: "svc-1"},
					SubnetIDs: *value.NewList([]value.Value[string]{
						value.New("subnet-1", 0, "", nil),
					}, 0, "", nil),
				},
			},
		},
		DirectConnect: directconnect.DirectConnect{
			GatewayAssociations: []directconnect.GatewayAssociation{
				{
					Resource:            resource.Resource{ID: "ga-1"},
					AssociatedGatewayID: value.New("tgw-1", 0, "", nil),
				},
			},
		},
	}

	a.PostProcess()

	mskTargets := append([]appautoscaling.Target(nil), a.MSK.Clusters[0].Relationships.AppAutoscalingTargets...)
	ddbTargets := append([]appautoscaling.Target(nil), a.DynamoDB.Tables[0].Relationships.AppAutoscalingTargets...)
	rgTargets := append([]*appautoscaling.Target(nil), a.ElastiCache.ReplicationGroups[0].Relationships.AppAutoscalingTargets...)
	eksLT := a.EKS.NodeGroups[0].Relationships.LaunchTemplate
	schTD := a.Scheduler.Schedules[0].Relationships.TaskDefinition
	cwTD := a.CloudWatch.EventTargets[0].Relationships.TaskDefinition
	pipeTD := a.Pipes.Pipes[0].Relationships.TaskDefinition
	gaTGW := a.DirectConnect.GatewayAssociations[0].Relationships.TransitGateway
	svcSubnets := append([]*ec2.Subnet(nil), a.ECS.Services[0].Relationships.Subnets...)

	a.PostProcess()

	assert.Equal(t, mskTargets, a.MSK.Clusters[0].Relationships.AppAutoscalingTargets)
	assert.Equal(t, ddbTargets, a.DynamoDB.Tables[0].Relationships.AppAutoscalingTargets)
	assert.Equal(t, rgTargets, a.ElastiCache.ReplicationGroups[0].Relationships.AppAutoscalingTargets)
	assert.Equal(t, eksLT, a.EKS.NodeGroups[0].Relationships.LaunchTemplate)
	assert.Equal(t, schTD, a.Scheduler.Schedules[0].Relationships.TaskDefinition)
	assert.Equal(t, cwTD, a.CloudWatch.EventTargets[0].Relationships.TaskDefinition)
	assert.Equal(t, pipeTD, a.Pipes.Pipes[0].Relationships.TaskDefinition)
	assert.Equal(t, gaTGW, a.DirectConnect.GatewayAssociations[0].Relationships.TransitGateway)
	assert.Equal(t, svcSubnets, a.ECS.Services[0].Relationships.Subnets)

	// Sanity-check we actually exercised the code paths.
	assert.Len(t, mskTargets, 1)
	assert.Len(t, ddbTargets, 1)
	assert.Len(t, rgTargets, 1)
	assert.NotNil(t, eksLT)
	assert.NotNil(t, schTD)
	assert.NotNil(t, cwTD)
	assert.NotNil(t, pipeTD)
	assert.NotNil(t, gaTGW)
	assert.Len(t, svcSubnets, 1)
}
