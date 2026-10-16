package dynamodb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type TableReplica struct {
	resource.Resource  `tree:"-"`
	TableClassOverride value.Value[TableClass] `tree:"table_class_override"`
}

type TableClass uint32

const (
	TableClassUnknown TableClass = iota
	TableClassStandard
	TableClassStandardInfrequentAccess
)
