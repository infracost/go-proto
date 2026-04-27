package aws

import (
	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/aws/apigateway"
	"github.com/infracost/go-proto/pkg/tree/aws/apigatewayv2"
	"github.com/infracost/go-proto/pkg/tree/aws/appautoscaling"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudfront"
	"github.com/infracost/go-proto/pkg/tree/aws/cloudformation"
	"github.com/infracost/go-proto/pkg/tree/aws/backup"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
	"github.com/infracost/go-proto/pkg/tree/aws/elasticbeanstalk"
 	"github.com/infracost/go-proto/pkg/tree/aws/elasticsearch"
	"github.com/infracost/go-proto/pkg/tree/aws/fsx"
)

type AWS struct {
	EC2                     ec2.EC2                        `tree:"ec2"`
	CertificateManager      acm.CertificateManager         `tree:"acm"`
	PCACertificateAuthority acmpca.PCACertificateAuthority `tree:"acmpca"`
	APIGateway              apigateway.APIGateway          `tree:"apigateway"`
	APIGatewayV2            apigatewayv2.APIGatewayV2      `tree:"apigatewayv2"`
	AppAutoScaling          appautoscaling.AppAutoScaling  `tree:"appautoscaling"`
	Backup                  backup.Backup                  `tree:"backup"`
	CloudFormation          cloudformation.CloudFormation  `tree:"cloudformation"`
  CloudFront              cloudfront.CloudFront          `tree:"cloudfront"`
	ElasticBeanstalk        elasticbeanstalk.ElasticBeanstalk `tree:"elasticbeanstalk"`
  Elasticsearch           elasticsearch.Elasticsearch    `tree:"elasticsearch"`
	FSx                     fsx.FSx                        `tree:"fsx"`
}

func (aws *AWS) PostProcess() {
	// acm
	// add the cert authorities to the certificate manager so that they can be linked to certificates
	aws.CertificateManager.AddCertificateAuthorities(&aws.PCACertificateAuthority)
}
