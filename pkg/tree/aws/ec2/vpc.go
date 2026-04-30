package ec2

import "github.com/infracost/go-proto/pkg/tree/resource"

type VPC struct {
	resource.Resource `tree:"-"`

	Relationships VPCRelationships `tree:"-"`
}

type VPCRelationships struct {
	VPCEndpoints []*VPCEndpoint
}
