package aws

import (
	"github.com/infracost/go-proto/pkg/tree/aws/ec2"
)

type AWS struct {
	EC2 ec2.EC2 `tree:"ec2"`
}
