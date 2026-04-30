package cosmosdb

type CosmosDB struct {
	Accounts         []Account         `tree:"accounts"`
	Databases        []Database        `tree:"databases"`
	CassandraTables  []CassandraTable  `tree:"cassandra_tables"`
	MongoCollections []MongoCollection `tree:"mongo_collections"`
}

func (s *CosmosDB) PostProcess() {
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
}
