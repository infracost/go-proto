package cosmosdb

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CassandraTable struct {
	resource.Resource `tree:"-"`
	Throughput        value.Int    `tree:"throughput"`
	MaxThroughput     value.Int    `tree:"max_throughput"`
	KeyspaceID        value.String `tree:"keyspace_id"`

	Relationships CassandraTableRelationships `tree:"-"`
}

type CassandraTableRelationships struct {
	Database *Database
}
