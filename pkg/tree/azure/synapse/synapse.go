package synapse

type Synapse struct {
	Workspaces []Workspace `tree:"workspaces"`
	SQLPools   []SQLPool   `tree:"sql_pools"`
	SparkPools []SparkPool `tree:"spark_pools"`
}

func (s *Synapse) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.SparkPools {
		s.SparkPools[i].Relationships.Workspace = nil
	}
	for i := range s.SQLPools {
		s.SQLPools[i].Relationships.Workspace = nil
	}

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
