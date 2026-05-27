package rds

type RDS struct {
	Clusters  []Cluster  `tree:"clusters"`
	Instances []Instance `tree:"instances"`
}

func (r *RDS) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range r.Instances {
		r.Instances[i].Relationships = InstanceRelationships{}
	}
	for i := range r.Clusters {
		r.Clusters[i].Relationships.Instances = nil
	}

	// link instances to clusters
	for i, instance := range r.Instances {
		if !instance.ClusterID.IsEmpty() {
			for j := range r.Clusters {
				if instance.ClusterID.Value() == r.Clusters[j].ID ||
					instance.ClusterID.Value() == r.Clusters[j].Identifier.Value() {
					r.Instances[i].Relationships.Cluster = &r.Clusters[j]
					r.Clusters[j].Relationships.Instances = append(r.Clusters[j].Relationships.Instances, &r.Instances[i])
					break
				}
			}
		}
	}

	// link replica source instances
	for i, instance := range r.Instances {
		if !instance.ReplicateSourceDBIdentifier.IsEmpty() {
			for j := range r.Instances {
				if instance.ReplicateSourceDBIdentifier.Value() == r.Instances[j].ID ||
					instance.ReplicateSourceDBIdentifier.Value() == r.Instances[j].Identifier.Value() {
					r.Instances[i].Relationships.ReplicateSourceDB = &r.Instances[j]
					break
				}
			}
		}
	}
}
