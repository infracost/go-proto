package inspector2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Enabler struct {
	resource.Resource `tree:"-"`
	EC2               value.Bool `tree:"ec2"`
	ECR               value.Bool `tree:"ecr"`
	Lambda            value.Bool `tree:"lambda"`
	LambdaCode        value.Bool `tree:"lambda_code"`
}
