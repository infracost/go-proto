package search

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SearchService struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	PartitionCount    value.Int    `tree:"partition_count"`
	ReplicaCount      value.Int    `tree:"replica_count"`
}
