package ecr

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LifecyclePolicy struct {
	resource.Resource `tree:"-"`
	RepositoryName    value.String `tree:"repository_name"`
}
