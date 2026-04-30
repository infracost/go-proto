package cosmosdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Database struct {
	resource.Resource `tree:"-"`
	Name              value.String `tree:"name"`
	Throughput        value.Int    `tree:"throughput"`
	MaxThroughput     value.Int    `tree:"max_throughput"`
	AccountName       value.String `tree:"account_name"`

	Relationships DatabaseRelationships `tree:"-"`
}

type DatabaseRelationships struct {
	Account *Account
}
