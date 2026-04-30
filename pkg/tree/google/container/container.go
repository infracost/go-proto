package container

type Container struct {
	Clusters   []Cluster   `tree:"clusters"`
	NodePools  []NodePool  `tree:"node_pools"`
	Registries []Registry  `tree:"registries"`
}

func (c *Container) PostProcess() {
	// Link NodePools to Clusters by ClusterID
	for i, np := range c.NodePools {
		if np.ClusterID.IsEmpty() {
			continue
		}
		for j := range c.Clusters {
			cluster := &c.Clusters[j]
			if np.ClusterID.Equal(cluster.ID) || np.ClusterID.Value() == cluster.Name.Value() {
				c.NodePools[i].Relationships.Cluster = cluster

				// Inherit Zones from Cluster if not set
				if np.Zones.IsEmpty() {
					c.NodePools[i].Zones = cluster.Zones
				}

				// Append to Cluster's NodePools
				cluster.Relationships.NodePools = append(cluster.Relationships.NodePools, &c.NodePools[i])
				break
			}
		}
	}
}
