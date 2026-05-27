package spanner

type Spanner struct {
	Instances []Instance `tree:"instances"`
	Databases []Database `tree:"databases"`
}

func (s *Spanner) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.Databases {
		s.Databases[i].Relationships.Instance = nil
	}

	// Link Databases to Instances by InstanceName
	for i, db := range s.Databases {
		if db.InstanceName.IsEmpty() {
			continue
		}
		for j := range s.Instances {
			inst := &s.Instances[j]
			if db.InstanceName.Equal(inst.ID) || db.InstanceName.Value() == inst.Name.Value() {
				s.Databases[i].Relationships.Instance = inst
				break
			}
		}
	}
}
