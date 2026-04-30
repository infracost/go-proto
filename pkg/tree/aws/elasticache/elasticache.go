package elasticache

type ElastiCache struct {
	Clusters          []Cluster          `tree:"clusters"`
	ReplicationGroups []ReplicationGroup `tree:"replication_groups"`
	ParameterGroups   []ParameterGroup   `tree:"parameter_groups"`
}

func (e *ElastiCache) PostProcess() {
	for i, cluster := range e.Clusters {
		// link cluster to parameter group
		for j, pg := range e.ParameterGroups {
			if pg.Name.Value() == cluster.ParameterGroupName.Value() || cluster.ParameterGroupName.Value() == pg.ID {
				e.Clusters[i].Relationships.ParameterGroup = &e.ParameterGroups[j]
				break
			}
		}

		// link cluster to replication group
		for j, rg := range e.ReplicationGroups {
			if rg.ID.Value() == cluster.ReplicationGroupID.Value() {
				e.Clusters[i].Relationships.ReplicationGroup = &e.ReplicationGroups[j]

				if cluster.Engine.IsDefaultOrEmpty() {
					e.Clusters[i].Engine = rg.Engine
				}
				if cluster.EngineVersion.IsDefaultOrEmpty() {
					e.Clusters[i].EngineVersion = rg.EngineVersion
				}
				if cluster.ParameterGroupName.IsDefaultOrEmpty() {
					e.Clusters[i].ParameterGroupName = rg.ParameterGroupName
				}

				// link to parameter group from replication group
				if e.Clusters[i].Relationships.ParameterGroup == nil {
					for k, pg := range e.ParameterGroups {
						if pg.Name.Value() == rg.ParameterGroupName.Value() || rg.ParameterGroupName.Value() == pg.ID {
							e.Clusters[i].Relationships.ParameterGroup = &e.ParameterGroups[k]
							break
						}
					}
				}
				break
			}
		}
	}
}
