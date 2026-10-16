package kendra

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DataSource struct {
	resource.Resource `tree:"-"`
	IndexID           value.String `tree:"index_id"`

	Relationships DataSourceRelationships `tree:"-"`
}

type DataSourceRelationships struct {
	Index *Index
	// Set on the first data source of each index, which carries per-index connector fees.
	FirstForIndex bool
}
