package msk

import (
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Cluster struct {
	resource.Resource      `tree:"-"`
	NumberOfBrokerNodes    value.Int    `tree:"number_of_broker_nodes"`
	BrokerNodeInstanceType value.String `tree:"broker_node_instance_type"`
	BrokerNodeEBSVolumeSizeGB value.Int `tree:"broker_node_ebs_volume_size"`

	Relationships ClusterRelationships `tree:"-"`
}

type ClusterRelationships struct {
	AppAutoscalingTargets []appautoscaling.Target
}
