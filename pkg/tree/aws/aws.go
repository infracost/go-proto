package aws

import (
	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/aws/apigateway"
	"github.com/infracost/go-proto/pkg/tree/aws/apigatewayv2"
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/aws/backup"
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
	"github.com/infracost/go-proto/pkg/tree/aws/fsx"
)

type AWS struct {
	EC2                     ec2.EC2                           `tree:"ec2"`
	CertificateManager      acm.CertificateManager            `tree:"acm"`
	PCACertificateAuthority acmpca.PCACertificateAuthority    `tree:"acmpca"`
	APIGateway              apigateway.APIGateway             `tree:"apigateway"`
	APIGatewayV2            apigatewayv2.APIGatewayV2         `tree:"apigatewayv2"`
	AppAutoScaling          appautoscaling.AppAutoScaling     `tree:"appautoscaling"`
	Backup                  backup.Backup                     `tree:"backup"`
	CloudFormation          cloudformation.CloudFormation     `tree:"cloudformation"`
	CloudFront              cloudfront.CloudFront             `tree:"cloudfront"`
	CloudHSMV2              cloudhsmv2.CloudHSMV2             `tree:"cloudhsmv2"`
	Cloudtrail              cloudtrail.CloudTrail             `tree:"cloudtrail"`
	CloudWatch              cloudwatch.CloudWatch             `tree:"cloudwatch"`
	CodeBuild               codebuild.CodeBuild               `tree:"codebuild"`
	Config                  config.Config                     `tree:"config"`
	DataTransfer            datatransfer.DataTransfer         `tree:"datatransfer"`
	DirectConnect           directconnect.DirectConnect       `tree:"directconnect"`
	DirectoryService        directoryservice.DirectoryService `tree:"directoryservice"`
	DMS                     dms.DMS                           `tree:"dms"`
	DocDB                   docdb.DocDB                       `tree:"docdb"`
	DynamoDB                dynamodb.DynamoDB                 `tree:"dynamodb"`
	ECR                     ecr.ECR                           `tree:"ecr"`
	ECS                     ecs.ECS                           `tree:"ecs"`
	EFS                     efs.EFS                           `tree:"efs"`
	EKS                     eks.EKS                           `tree:"eks"`
	ElastiCache             elasticache.ElastiCache           `tree:"elasticache"`
	ElasticBeanstalk        elasticbeanstalk.ElasticBeanstalk `tree:"elasticbeanstalk"`
	Elasticsearch           elasticsearch.Elasticsearch       `tree:"elasticsearch"`
	FSx                     fsx.FSx                           `tree:"fsx"`
}

func (aws *AWS) PostProcess() {
	// acm
	// add the cert authorities to the certificate manager so that they can be linked to certificates
	aws.CertificateManager.AddCertificateAuthorities(&aws.PCACertificateAuthority)

	// elasticache
	// link parameter groups and replication groups to clusters
	aws.ElastiCache.PostProcess()
	// ecs
	// link capacity providers, task definitions, and clusters to services
	aws.ECS.PostProcess()
	// ecr
	// link lifecycle policies to repositories
	aws.ECR.PostProcess()
}
