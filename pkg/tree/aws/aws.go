package aws

import (
	"github.com/infracost/go-proto/pkg/tree/aws/acm"
	"github.com/infracost/go-proto/pkg/tree/aws/acmpca"
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
)

type AWS struct {
	EC2                     ec2.EC2                        `tree:"ec2"`
	CertificateManager      acm.CertificateManager         `tree:"acm"`
	PCACertificateAuthority acmpca.PCACertificateAuthority `tree:"acmpca"`
}

func (aws *AWS) PostProcess() {
	// acm
	// add the cert authorities to the certificate manager so that they can be linked to certificates
	aws.CertificateManager.AddCertificateAuthorities(&aws.PCACertificateAuthority)
}
