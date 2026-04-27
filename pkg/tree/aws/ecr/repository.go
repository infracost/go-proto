package ecr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Repository struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`

	Relationships RepositoryRelationships `tree:"-"`
}

type RepositoryRelationships struct {
	LifecyclePolicies []*LifecyclePolicy
}
