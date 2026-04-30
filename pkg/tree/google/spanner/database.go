package spanner

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Database struct {
	resource.Resource `tree:"-"`
	InstanceName      value.String `tree:"instance_name"`

	Relationships DatabaseRelationships `tree:"-"`
}

type DatabaseRelationships struct {
	Instance *Instance
}
