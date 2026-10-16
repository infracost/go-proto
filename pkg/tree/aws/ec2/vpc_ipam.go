package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type VPCIPAM struct {
	resource.Resource `tree:"-"`
	Tier              value.String `tree:"tier"`
}
