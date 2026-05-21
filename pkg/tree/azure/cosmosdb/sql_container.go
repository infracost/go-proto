package cosmosdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type SQLContainer struct {
	resource.Resource `tree:"-"`
	Throughput        value.Int    `tree:"throughput"`
	MaxThroughput     value.Int    `tree:"max_throughput"`
	DatabaseName      value.String `tree:"database_name"`
	AccountName       value.String `tree:"account_name"`

	Relationships SQLContainerRelationships `tree:"-"`
}

type SQLContainerRelationships struct {
	Database *Database
	Account  *Account
}
