package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type ElasticIPAssociation struct {
	resource.Resource `tree:"-"`
	AllocationID      value.String `tree:"allocation_id"`

	Relationships ElasticIPAssociationRelationships `tree:"-"`
}

type ElasticIPAssociationRelationships struct {
	ElasticIP *ElasticIP
}
