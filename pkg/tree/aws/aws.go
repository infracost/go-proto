package aws

import (
	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/aws/apigateway"
	"github.com/infracost/go-proto/pkg/tree/aws/apigatewayv2"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
)

type AWS struct {
	EC2                     ec2.EC2                        `tree:"ec2"`
	CertificateManager      acm.CertificateManager         `tree:"acm"`
	PCACertificateAuthority acmpca.PCACertificateAuthority `tree:"acmpca"`
	APIGateway              apigateway.API                 `tree:"apigateway"`
	APIGatewayV2            apigatewayv2.API               `tree:"apigatewayv2"`
}

func (aws *AWS) PostProcess() {
	// acm
	// add the cert authorities to the certificate manager so that they can be linked to certificates
	aws.CertificateManager.AddCertificateAuthorities(&aws.PCACertificateAuthority)
}
