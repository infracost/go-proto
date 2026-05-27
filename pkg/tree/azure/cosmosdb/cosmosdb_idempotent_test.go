package cosmosdb

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	s := &CosmosDB{
		Accounts: []Account{
			{
				Resource: resource.Resource{ID: "acc-1"},
				Name:     value.New("acc-1", 0, "", nil),
			},
		},
		Databases: []Database{
			{
				Resource:    resource.Resource{ID: "db-1"},
				Name:        value.New("db-1", 0, "", nil),
				AccountName: value.New("acc-1", 0, "", nil),
			},
		},
		CassandraTables: []CassandraTable{
			{
				Resource:   resource.Resource{ID: "ct-1"},
				KeyspaceID: value.New("db-1", 0, "", nil),
			},
		},
		MongoCollections: []MongoCollection{
			{
				Resource:     resource.Resource{ID: "mc-1"},
				DatabaseName: value.New("db-1", 0, "", nil),
			},
		},
		GremlinGraphs: []GremlinGraph{
			{
				Resource:     resource.Resource{ID: "gg-1"},
				DatabaseName: value.New("db-1", 0, "", nil),
				AccountName:  value.New("acc-1", 0, "", nil),
			},
		},
		SQLContainers: []SQLContainer{
			{
				Resource:     resource.Resource{ID: "sc-1"},
				DatabaseName: value.New("db-1", 0, "", nil),
				AccountName:  value.New("acc-1", 0, "", nil),
			},
		},
	}

	s.PostProcess()
	dbAcc := s.Databases[0].Relationships.Account
	ctDB := s.CassandraTables[0].Relationships.Database
	mcDB := s.MongoCollections[0].Relationships.Database
	ggDB := s.GremlinGraphs[0].Relationships.Database
	ggAcc := s.GremlinGraphs[0].Relationships.Account
	scDB := s.SQLContainers[0].Relationships.Database
	scAcc := s.SQLContainers[0].Relationships.Account

	s.PostProcess()
	assert.Equal(t, dbAcc, s.Databases[0].Relationships.Account)
	assert.Equal(t, ctDB, s.CassandraTables[0].Relationships.Database)
	assert.Equal(t, mcDB, s.MongoCollections[0].Relationships.Database)
	assert.Equal(t, ggDB, s.GremlinGraphs[0].Relationships.Database)
	assert.Equal(t, ggAcc, s.GremlinGraphs[0].Relationships.Account)
	assert.Equal(t, scDB, s.SQLContainers[0].Relationships.Database)
	assert.Equal(t, scAcc, s.SQLContainers[0].Relationships.Account)

	assert.NotNil(t, dbAcc)
	assert.NotNil(t, ctDB)
	assert.NotNil(t, mcDB)
	assert.NotNil(t, ggDB)
	assert.NotNil(t, ggAcc)
	assert.NotNil(t, scDB)
	assert.NotNil(t, scAcc)
}
