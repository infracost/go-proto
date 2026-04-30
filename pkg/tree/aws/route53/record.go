package route53

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Record struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	AliasName         value.String `tree:"alias_name"`

	Relationships RecordRelationships `tree:"-"`
}

type RecordRelationships struct {
	AliasRecord *Record
}
