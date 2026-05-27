package neptune

type Neptune struct {
	Clusters         []Cluster         `tree:"clusters"`
	ClusterInstances []ClusterInstance `tree:"cluster_instances"`
	ClusterSnapshots []ClusterSnapshot `tree:"cluster_snapshots"`
}

func (n *Neptune) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range n.ClusterInstances {
		n.ClusterInstances[i].Relationships.Cluster = nil
	}
	for i := range n.ClusterSnapshots {
		n.ClusterSnapshots[i].Relationships.Cluster = nil
	}

	// link cluster instances to clusters
	for i, instance := range n.ClusterInstances {
		for j := range n.Clusters {
			if instance.ClusterIdentifier.Value() == n.Clusters[j].Identifier.Value() ||
				instance.ClusterIdentifier.Value() == n.Clusters[j].ID {
				n.ClusterInstances[i].Relationships.Cluster = &n.Clusters[j]
				break
			}
		}
	}

	// link cluster snapshots to clusters
	for i, snapshot := range n.ClusterSnapshots {
		for j := range n.Clusters {
			if snapshot.DBClusterIdentifier.Value() == n.Clusters[j].Identifier.Value() ||
				snapshot.DBClusterIdentifier.Value() == n.Clusters[j].ID {
				n.ClusterSnapshots[i].Relationships.Cluster = &n.Clusters[j]
				break
			}
		}
	}
}
