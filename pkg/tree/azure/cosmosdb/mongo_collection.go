package cosmosdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MongoCollection struct {
	resource.Resource `tree:"-"`
	Throughput        value.Int    `tree:"throughput"`
	MaxThroughput     value.Int    `tree:"max_throughput"`
	DatabaseName      value.String `tree:"database_name"`

	Relationships MongoCollectionRelationships `tree:"-"`
}

type MongoCollectionRelationships struct {
	Database *Database
}
