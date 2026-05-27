package cosmosdb

type CosmosDB struct {
	Accounts         []Account         `tree:"accounts"`
	Databases        []Database        `tree:"databases"`
	CassandraTables  []CassandraTable  `tree:"cassandra_tables"`
	MongoCollections []MongoCollection `tree:"mongo_collections"`
	GremlinGraphs    []GremlinGraph    `tree:"gremlin_graphs"`
	SQLContainers    []SQLContainer    `tree:"sql_containers"`
}

func (s *CosmosDB) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.Databases {
		s.Databases[i].Relationships.Account = nil
	}
	for i := range s.CassandraTables {
		s.CassandraTables[i].Relationships.Database = nil
	}
	for i := range s.MongoCollections {
		s.MongoCollections[i].Relationships.Database = nil
	}
	for i := range s.GremlinGraphs {
		s.GremlinGraphs[i].Relationships = GremlinGraphRelationships{}
	}
	for i := range s.SQLContainers {
		s.SQLContainers[i].Relationships = SQLContainerRelationships{}
	}

	for i, db := range s.Databases {
		for j := range s.Accounts {
			if db.AccountName.Value() == s.Accounts[j].Name.Value() {
				s.Databases[i].Relationships.Account = &s.Accounts[j]
				break
			}
		}
	}

	for i, ct := range s.CassandraTables {
		for j := range s.Databases {
			if ct.KeyspaceID.Value() == s.Databases[j].ID {
				s.CassandraTables[i].Relationships.Database = &s.Databases[j]
				break
			}
		}
	}

	for i, mc := range s.MongoCollections {
		for j := range s.Databases {
			if mc.DatabaseName.Value() == s.Databases[j].Name.Value() {
				s.MongoCollections[i].Relationships.Database = &s.Databases[j]
				break
			}
		}
	}

	for i, g := range s.GremlinGraphs {
		for j := range s.Databases {
			if g.DatabaseName.Value() == s.Databases[j].Name.Value() {
				s.GremlinGraphs[i].Relationships.Database = &s.Databases[j]
				break
			}
		}
		for j := range s.Accounts {
			if g.AccountName.Value() == s.Accounts[j].Name.Value() {
				s.GremlinGraphs[i].Relationships.Account = &s.Accounts[j]
				break
			}
		}
	}

	for i, c := range s.SQLContainers {
		for j := range s.Databases {
			if c.DatabaseName.Value() == s.Databases[j].Name.Value() {
				s.SQLContainers[i].Relationships.Database = &s.Databases[j]
				break
			}
		}
		for j := range s.Accounts {
			if c.AccountName.Value() == s.Accounts[j].Name.Value() {
				s.SQLContainers[i].Relationships.Account = &s.Accounts[j]
				break
			}
		}
	}
}
