package aws

import (
	"fmt"

	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/aws/apigateway"
	"github.com/infracost/go-proto/pkg/tree/aws/apigatewayv2"
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/aws/backup"
	"github.com/infracost/go-proto/pkg/tree/aws/batch"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudformation"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudfront"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudhsmv2"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudtrail"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudwatch"
	"github.com/infracost/go-proto/pkg/tree/aws/codebuild"
	"github.com/infracost/go-proto/pkg/tree/aws/config"
	"github.com/infracost/go-proto/pkg/tree/aws/datatransfer"
	"github.com/infracost/go-proto/pkg/tree/aws/directconnect"
	"github.com/infracost/go-proto/pkg/tree/aws/directoryservice"
	"github.com/infracost/go-proto/pkg/tree/aws/dms"
	"github.com/infracost/go-proto/pkg/tree/aws/docdb"
	"github.com/infracost/go-proto/pkg/tree/aws/dynamodb"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/aws/ecr"
	"github.com/infracost/go-proto/pkg/tree/aws/ecs"
	"github.com/infracost/go-proto/pkg/tree/aws/efs"
	"github.com/infracost/go-proto/pkg/tree/aws/eks"
	"github.com/infracost/go-proto/pkg/tree/aws/elasticache"
	"github.com/infracost/go-proto/pkg/tree/aws/elasticbeanstalk"
	"github.com/infracost/go-proto/pkg/tree/aws/elasticsearch"
	"github.com/infracost/go-proto/pkg/tree/aws/emr"
	"github.com/infracost/go-proto/pkg/tree/aws/fsx"
	"github.com/infracost/go-proto/pkg/tree/aws/globalaccelerator"
	"github.com/infracost/go-proto/pkg/tree/aws/glue"
	"github.com/infracost/go-proto/pkg/tree/aws/kinesis"
	"github.com/infracost/go-proto/pkg/tree/aws/kinesisanalytics"
	"github.com/infracost/go-proto/pkg/tree/aws/kinesisanalyticsv2"
	"github.com/infracost/go-proto/pkg/tree/aws/kms"
	"github.com/infracost/go-proto/pkg/tree/aws/lambda"
	"github.com/infracost/go-proto/pkg/tree/aws/lightsail"
	"github.com/infracost/go-proto/pkg/tree/aws/mq"
	"github.com/infracost/go-proto/pkg/tree/aws/msk"
	"github.com/infracost/go-proto/pkg/tree/aws/mwaa"
	"github.com/infracost/go-proto/pkg/tree/aws/neptune"
	"github.com/infracost/go-proto/pkg/tree/aws/networkfirewall"
	"github.com/infracost/go-proto/pkg/tree/aws/pipes"
	"github.com/infracost/go-proto/pkg/tree/aws/rds"
	"github.com/infracost/go-proto/pkg/tree/aws/redshift"
	"github.com/infracost/go-proto/pkg/tree/aws/route53"
	"github.com/infracost/go-proto/pkg/tree/aws/s3"
	"github.com/infracost/go-proto/pkg/tree/aws/sagemaker"
	"github.com/infracost/go-proto/pkg/tree/aws/scheduler"
	"github.com/infracost/go-proto/pkg/tree/aws/secretsmanager"
	"github.com/infracost/go-proto/pkg/tree/aws/sns"
	"github.com/infracost/go-proto/pkg/tree/aws/sqs"
	"github.com/infracost/go-proto/pkg/tree/aws/ssm"
	"github.com/infracost/go-proto/pkg/tree/aws/stepfunctions"
	"github.com/infracost/go-proto/pkg/tree/aws/transfer"
	"github.com/infracost/go-proto/pkg/tree/aws/waf"
)

type AWS struct {
	EC2                     ec2.EC2                               `tree:"ec2"`
	CertificateManager      acm.CertificateManager                `tree:"acm"`
	PCACertificateAuthority acmpca.PCACertificateAuthority        `tree:"acmpca"`
	APIGateway              apigateway.APIGateway                 `tree:"apigateway"`
	APIGatewayV2            apigatewayv2.APIGatewayV2             `tree:"apigatewayv2"`
	AppAutoScaling          appautoscaling.AppAutoScaling         `tree:"appautoscaling"`
	Backup                  backup.Backup                         `tree:"backup"`
	Batch                   batch.Batch                           `tree:"batch"`
	CloudFormation          cloudformation.CloudFormation         `tree:"cloudformation"`
	CloudFront              cloudfront.CloudFront                 `tree:"cloudfront"`
	CloudHSMV2              cloudhsmv2.CloudHSMV2                 `tree:"cloudhsmv2"`
	Cloudtrail              cloudtrail.CloudTrail                 `tree:"cloudtrail"`
	CloudWatch              cloudwatch.CloudWatch                 `tree:"cloudwatch"`
	CodeBuild               codebuild.CodeBuild                   `tree:"codebuild"`
	Config                  config.Config                         `tree:"config"`
	DataTransfer            datatransfer.DataTransfer             `tree:"datatransfer"`
	DirectConnect           directconnect.DirectConnect           `tree:"directconnect"`
	DirectoryService        directoryservice.DirectoryService     `tree:"directoryservice"`
	DMS                     dms.DMS                               `tree:"dms"`
	DocDB                   docdb.DocDB                           `tree:"docdb"`
	DynamoDB                dynamodb.DynamoDB                     `tree:"dynamodb"`
	ECR                     ecr.ECR                               `tree:"ecr"`
	ECS                     ecs.ECS                               `tree:"ecs"`
	EFS                     efs.EFS                               `tree:"efs"`
	EKS                     eks.EKS                               `tree:"eks"`
	ElastiCache             elasticache.ElastiCache               `tree:"elasticache"`
	ElasticBeanstalk        elasticbeanstalk.ElasticBeanstalk     `tree:"elasticbeanstalk"`
	Elasticsearch           elasticsearch.Elasticsearch           `tree:"elasticsearch"`
	EMR                     emr.EMR                               `tree:"emr"`
	FSx                     fsx.FSx                               `tree:"fsx"`
	GlobalAccelerator       globalaccelerator.GlobalAccelerator   `tree:"globalaccelerator"`
	Glue                    glue.Glue                             `tree:"glue"`
	Kinesis                 kinesis.Kinesis                       `tree:"kinesis"`
	KinesisAnalytics        kinesisanalytics.KinesisAnalytics     `tree:"kinesisanalytics"`
	KinesisAnalyticsV2      kinesisanalyticsv2.KinesisAnalyticsV2 `tree:"kinesisanalyticsv2"`
	KMS                     kms.KMS                               `tree:"kms"`
	Lambda                  lambda.Lambda                         `tree:"lambda"`
	Lightsail               lightsail.Lightsail                   `tree:"lightsail"`
	MQ                      mq.MQ                                 `tree:"mq"`
	MSK                     msk.MSK                               `tree:"msk"`
	MWAA                    mwaa.MWAA                             `tree:"mwaa"`
	Neptune                 neptune.Neptune                       `tree:"neptune"`
	NetworkFirewall         networkfirewall.NetworkFirewall       `tree:"networkfirewall"`
	Pipes                   pipes.Pipes                           `tree:"pipes"`
	RDS                     rds.RDS                               `tree:"rds"`
	Redshift                redshift.Redshift                     `tree:"redshift"`
	Route53                 route53.Route53                       `tree:"route53"`
	S3                      s3.S3                                 `tree:"s3"`
	SageMaker               sagemaker.SageMaker                   `tree:"sagemaker"`
	Scheduler               scheduler.Scheduler                   `tree:"scheduler"`
	SecretsManager          secretsmanager.SecretsManager         `tree:"secretsmanager"`
	SNS                     sns.SNS                               `tree:"sns"`
	SQS                     sqs.SQS                               `tree:"sqs"`
	SSM                     ssm.SSM                               `tree:"ssm"`
	StepFunctions           stepfunctions.StepFunctions           `tree:"stepfunctions"`
	Transfer                transfer.Transfer                     `tree:"transfer"`
	WAF                     waf.WAF                               `tree:"waf"`
}

func (aws *AWS) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range aws.MSK.Clusters {
		aws.MSK.Clusters[i].Relationships.AppAutoscalingTargets = nil
	}
	for i := range aws.DynamoDB.Tables {
		aws.DynamoDB.Tables[i].Relationships.AppAutoscalingTargets = nil
	}
	for i := range aws.ElastiCache.ReplicationGroups {
		aws.ElastiCache.ReplicationGroups[i].Relationships.AppAutoscalingTargets = nil
	}
	for i := range aws.EKS.NodeGroups {
		aws.EKS.NodeGroups[i].Relationships.LaunchTemplate = nil
	}
	for i := range aws.Scheduler.Schedules {
		aws.Scheduler.Schedules[i].Relationships.TaskDefinition = nil
	}
	for i := range aws.CloudWatch.EventTargets {
		aws.CloudWatch.EventTargets[i].Relationships.TaskDefinition = nil
	}
	for i := range aws.Pipes.Pipes {
		aws.Pipes.Pipes[i].Relationships.TaskDefinition = nil
	}
	for i := range aws.DirectConnect.GatewayAssociations {
		aws.DirectConnect.GatewayAssociations[i].Relationships.TransitGateway = nil
	}
	for i := range aws.ECS.Services {
		aws.ECS.Services[i].Relationships.Subnets = nil
	}

	// acm - link cert authorities to certificate manager
	aws.CertificateManager.AddCertificateAuthorities(&aws.PCACertificateAuthority)

	// NOTE: Service-level PostProcess() methods (EC2, ECS, RDS, S3, etc.) are NOT called
	// here — they are invoked automatically by the reflective tree walker in tree.go.
	// Only cross-service relationship wiring belongs in this method.

	// cross-service: link app autoscaling targets to MSK clusters
	for i, cluster := range aws.MSK.Clusters {
		for _, target := range aws.AppAutoScaling.Targets {
			if target.ResourceID.Value() == cluster.ID {
				aws.MSK.Clusters[i].Relationships.AppAutoscalingTargets = append(
					aws.MSK.Clusters[i].Relationships.AppAutoscalingTargets, target,
				)
			}
		}
	}

	// cross-service: link app autoscaling targets to DynamoDB tables
	for i, table := range aws.DynamoDB.Tables {
		for _, target := range aws.AppAutoScaling.Targets {
			if target.ResourceID.Value() == fmt.Sprintf("table/%s", table.Name.Value()) {
				aws.DynamoDB.Tables[i].Relationships.AppAutoscalingTargets = append(
					aws.DynamoDB.Tables[i].Relationships.AppAutoscalingTargets, target,
				)
			}
		}
	}

	// cross-service: link app autoscaling targets to ElastiCache replication groups
	for i, rg := range aws.ElastiCache.ReplicationGroups {
		for j := range aws.AppAutoScaling.Targets {
			target := &aws.AppAutoScaling.Targets[j]
			if target.ResourceID.Value() == fmt.Sprintf("replication-group/%s", rg.ID.Value()) ||
				target.ResourceID.Value() == fmt.Sprintf("replication-group/%s", rg.Resource.ID) {
				aws.ElastiCache.ReplicationGroups[i].Relationships.AppAutoscalingTargets = append(
					aws.ElastiCache.ReplicationGroups[i].Relationships.AppAutoscalingTargets, target,
				)
			}
		}
	}

	// cross-service: link EKS node groups to EC2 launch templates
	for i, nodeGroup := range aws.EKS.NodeGroups {
		for j := range aws.EC2.LaunchTemplates {
			lt := &aws.EC2.LaunchTemplates[j]
			if (!nodeGroup.LaunchTemplateID.IsEmpty() && nodeGroup.LaunchTemplateID.Value() == lt.ID) ||
				(!nodeGroup.LaunchTemplateName.IsEmpty() && nodeGroup.LaunchTemplateName.Value() == lt.Name.Value()) ||
				(!nodeGroup.LaunchTemplateID.IsEmpty() && nodeGroup.LaunchTemplateID.Value() == lt.Name.Value()) {
				aws.EKS.NodeGroups[i].Relationships.LaunchTemplate = lt
				break
			}
		}
	}

	// cross-service: link scheduler schedules to ECS task definitions
	for i, schedule := range aws.Scheduler.Schedules {
		for j := range aws.ECS.TaskDefinitions {
			if aws.ECS.TaskDefinitions[j].ID == schedule.TaskDefinitionARN.Value() {
				aws.Scheduler.Schedules[i].Relationships.TaskDefinition = &aws.ECS.TaskDefinitions[j]
				break
			}
		}
	}

	// cross-service: link cloudwatch event targets to ECS task definitions
	for i, eventTarget := range aws.CloudWatch.EventTargets {
		for j := range aws.ECS.TaskDefinitions {
			if eventTarget.TaskDefinitionID.Value() == aws.ECS.TaskDefinitions[j].ID {
				aws.CloudWatch.EventTargets[i].Relationships.TaskDefinition = &aws.ECS.TaskDefinitions[j]
				break
			}
		}
	}

	// cross-service: link pipes to ECS task definitions
	for i, pipe := range aws.Pipes.Pipes {
		for j := range aws.ECS.TaskDefinitions {
			if pipe.TaskDefinitionID.Value() == aws.ECS.TaskDefinitions[j].ID {
				aws.Pipes.Pipes[i].Relationships.TaskDefinition = &aws.ECS.TaskDefinitions[j]
				break
			}
		}
	}

	// cross-service: link direct connect gateway associations to EC2 transit gateways
	for i, association := range aws.DirectConnect.GatewayAssociations {
		for j := range aws.EC2.TransitGateways {
			if aws.EC2.TransitGateways[j].ID == association.AssociatedGatewayID.Value() {
				aws.DirectConnect.GatewayAssociations[i].Relationships.TransitGateway = &aws.EC2.TransitGateways[j]
				break
			}
		}
	}

	// cross-service: link ECS services to EC2 subnets
	for i, svc := range aws.ECS.Services {
		for j := range aws.EC2.Subnets {
			if svc.SubnetIDs.Contains(aws.EC2.Subnets[j].ID) {
				aws.ECS.Services[i].Relationships.Subnets = append(
					aws.ECS.Services[i].Relationships.Subnets, &aws.EC2.Subnets[j],
				)
			}
		}
	}
}
