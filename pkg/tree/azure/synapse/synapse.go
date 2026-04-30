package synapse

type Synapse struct {
	Workspaces []Workspace `tree:"workspaces"`
	SQLPools   []SQLPool   `tree:"sql_pools"`
	SparkPools []SparkPool `tree:"spark_pools"`
}

func (s *Synapse) PostProcess() {
	for i, pool := range s.SparkPools {
		for j := range s.Workspaces {
			if pool.WorkspaceID.Value() == s.Workspaces[j].ID {
				s.SparkPools[i].Relationships.Workspace = &s.Workspaces[j]
				break
			}
		}
	}

	for i, pool := range s.SQLPools {
		for j := range s.Workspaces {
			if pool.WorkspaceID.Value() == s.Workspaces[j].ID {
				s.SQLPools[i].Relationships.Workspace = &s.Workspaces[j]
				break
			}
		}
	}
}
